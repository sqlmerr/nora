package middleware

import (
	"nora/internal/config"

	e "nora/internal/error"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(config *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.JwtSecret)},
		ErrorHandler: jwtError,
		ContextKey:   "jwtToken",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(e.New(err.Error(), fiber.StatusBadRequest))
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(e.ErrUnauthorized)
}
