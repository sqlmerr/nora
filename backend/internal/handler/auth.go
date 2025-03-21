package handler

import (
	"errors"
	"fmt"
	"net/http"
	"nora/internal/config"
	e "nora/internal/error"
	"nora/internal/model"
	"nora/internal/service"

	"github.com/gofiber/fiber/v2"
)

func registerUser(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	var data model.UserCreate
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	user, err := s.Register(c.UserContext(), &data)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.Status(201).JSON(user)
}

func login(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	var data model.UserLogin
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	token, err := s.Login(c.UserContext(), &data)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(fiber.Map{"token": token})
}

func getMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)
	return c.JSON(user)
}

func connectTelegram(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	cfg := c.Locals("config").(*config.Config)
	secretToken := c.Get("secret-token")

	if secretToken != cfg.SecretToken {
		fmt.Println(secretToken, cfg.SecretToken)
		return c.Status(http.StatusUnauthorized).JSON(e.ErrInvalidToken)
	}

	var data model.UserTelegramConnect
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	err := s.ConnectTelegram(c.UserContext(), data.UserID, data.TelegramID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(fiber.Map{"message": "telegram connected", "ok": true})
}
