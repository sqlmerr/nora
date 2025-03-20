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

func listTasks(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)

	id := c.Params("groupId")
	groupId, err := uuid.Parse(id)
	if err != nil || groupId == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	tasks, err := s.FindAllTasksByGroup(c.Context(), groupId, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}
	return c.JSON(tasks)
}

func getTask(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
	if err != nil || taskID == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	task, err := s.FindOneTask(c.Context(), taskID, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(task)
}

func createTask(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	var data model.TaskCreate
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	task, err := s.CreateTask(c.Context(), &data, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.Status(201).JSON(task)
}

func deleteTask(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
	if err != nil || taskID == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	err = s.DeleteTask(c.Context(), taskID, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(fiber.Map{"message": "task deleted", "ok": true})
}

func updateTask(c *fiber.Ctx) error {
	s := c.Locals("service").(*service.Service)
	user := c.Locals("user").(*model.User)
	id := c.Params("id")
	taskID, err := uuid.Parse(id)
	if err != nil || taskID == uuid.Nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	var data model.TaskUpdate
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(e.ErrInvalidBody)
	}

	err = s.UpdateTask(c.Context(), taskID, &data, user.ID)
	if err != nil {
		var apiErr e.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status).JSON(err)
		}
		return c.Status(http.StatusInternalServerError).JSON(e.New(err.Error(), 500))
	}

	return c.JSON(fiber.Map{"message": "task updated", "ok": true})
}
