package handler

import (
	"errors"
	"net/http"
	e "nora/internal/error"
	"nora/internal/model"
	"nora/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func createSpace(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)
	s := c.Locals("service").(*service.Service)
	var data model.SpaceCreate
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	data.UserID = user.ID

	space, err := s.CreateSpace(c.Context(), &data)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.Status(201).JSON(space)
}

func listSpaces(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)

	spaces, err := s.ListSpaces(c.Context(), user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(spaces)
}

func getSpace(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	id := c.Params("id")
	spaceID, err := uuid.Parse(id)
	if err != nil || spaceID == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	space, err := s.GetSpace(c.Context(), spaceID, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(space)
}

func deleteSpace(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	id := c.Params("id")
	spaceID, err := uuid.Parse(id)
	if err != nil || spaceID == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	err = s.DeleteSpace(c.Context(), spaceID, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(fiber.Map{"message": "space deleted", "ok": true})
}
