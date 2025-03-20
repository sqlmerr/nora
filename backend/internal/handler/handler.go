package handler

import (
	"nora/internal/config"
	mw "nora/internal/middleware"
	"nora/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func root(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*zap.Logger)
	logger.Info("hello")
	return c.SendString("Ok")
}

func New(logger *zap.Logger, s *service.Service, config *config.Config) *fiber.App {
	jwtMiddleware := mw.JwtMiddleware(config)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logger)
		c.Locals("service", s)
		c.Locals("config", config)
		return c.Next()
	})
	api := app.Group("/api")
	api.Get("/", root)

	auth := api.Group("/auth")
	auth.Post("/register", registerUser)
	auth.Post("/login", login)
	auth.Get("/me", jwtMiddleware, mw.UserGetter, getMe)

	spaces := api.Group("/spaces", jwtMiddleware, mw.UserGetter)
	spaces.Get("/", listSpaces)
	spaces.Get("/:id", getSpace)
	spaces.Post("/", createSpace)
	spaces.Delete("/:id", deleteSpace)

	spaces.Get("/:spaceId/groups", listTaskGroups)

	tasks := api.Group("/tasks", jwtMiddleware, mw.UserGetter)
	tasks.Get("/:groupId", listTasks)
	tasks.Post("/", createTask)
	tasks.Get("/:id", getTask)
	tasks.Put("/:id", updateTask)
	tasks.Delete("/:id", deleteTask)

	groups := api.Group("/groups", jwtMiddleware, mw.UserGetter)
	groups.Get("/:spaceId", listTaskGroups)
	groups.Post("/", createTaskGroup)

	return app
}
