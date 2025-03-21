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

func listTaskGroups(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)

	id := c.Params("spaceId")
	spaceId, err := uuid.Parse(id)
	if err != nil || spaceId == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	groups, err := s.FindAllTaskGroupsBySpace(c.UserContext(), spaceId, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}
	return c.JSON(groups)
}

func createTaskGroup(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	var data model.TaskGroupCreate
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	group, err := s.CreateTaskGroup(c.UserContext(), &data, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.Status(201).JSON(group)
}
