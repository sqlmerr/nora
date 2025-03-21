package handler

import (
	"errors"
	"net/http"
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
