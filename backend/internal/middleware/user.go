package middleware

import (
	"net/http"
	"nora/internal/service"

	e "nora/internal/error"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserGetter(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)

	token := c.Locals("jwtToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)

	user, err := s.FindOneUserByUsername(c.Context(), username)
	if err != nil || user == nil {
		return c.Status(http.StatusForbidden).JSON(e.ErrInvalidToken)
	}

	c.Locals("user", user)

	return c.Next()
}
