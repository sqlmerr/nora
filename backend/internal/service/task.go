package service

import (
	"context"
	e "nora/internal/error"
	"nora/internal/model"

	"github.com/google/uuid"
)

func (s *Service) CreateTask(ctx context.Context, data *model.TaskCreate, userID uuid.UUID) (*model.Task, error) {
	if data.GroupID == uuid.Nil {
		return nil, e.ErrInvalidBody
	}
	group, err := s.FindOneTaskGroup(ctx, data.GroupID, userID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, e.ErrNotFound
	}

	id, err := s.t.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return s.t.FindOne(ctx, *id)
}

func (s *Service) FindOneTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Task, error) {
	task, err := s.t.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	group, err := s.FindOneTaskGroup(ctx, task.GroupID, userID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, e.ErrNotFound
	}
	return task, nil
}

func (s *Service) FindAllTasksByGroup(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) ([]model.Task, error) {
	group, err := s.FindOneTaskGroup(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, e.ErrNotFound
	}

	return s.t.FindAllByGroup(ctx, groupID)
}

func (s *Service) UpdateTask(ctx context.Context, id uuid.UUID, data *model.TaskUpdate, userID uuid.UUID) error {
	task, err := s.t.FindOne(ctx, id)
	if err != nil || task == nil {
		return err
	}
	group, err := s.FindOneTaskGroup(ctx, task.GroupID, userID)
	if err != nil {
		return err
	}
	space, err := s.GetSpace(ctx, group.SpaceID, userID)
	if err != nil {
		return err
	}
	if space.UserID != userID {
		return e.ErrForbidden
	}
	if data.GroupID != uuid.Nil && data.GroupID != group.ID {
		newGroup, err := s.FindOneTaskGroup(ctx, data.GroupID, userID)
		if err != nil {
			return err
		}
		if newGroup == nil {
			return e.ErrNotFound
		}
	}

	return s.t.Update(ctx, id, data)
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	task, err := s.t.FindOne(ctx, id)
	if err != nil || task == nil {
		return err
	}

	return s.t.Delete(ctx, id)
}
