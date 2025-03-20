package service

import (
	"context"
	e "nora/internal/error"
	"nora/internal/model"

	"github.com/google/uuid"
)

func (s *Service) CreateTaskGroup(ctx context.Context, data *model.TaskGroupCreate, userID uuid.UUID) (*model.TaskGroup, error) {
	space, err := s.sp.FindOne(ctx, data.SpaceID)
	if err != nil {
		return nil, err
	}
	if space.UserID != userID {
		return nil, e.ErrForbidden
	}

	id, err := s.tg.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return s.tg.FindOne(ctx, *id)
}

func (s *Service) FindOneTaskGroup(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.TaskGroup, error) {

	taskGroup, err := s.tg.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if taskGroup == nil {
		return nil, e.ErrNotFound
	}
	space, err := s.GetSpace(ctx, taskGroup.SpaceID, userID)
	if err != nil {
		return nil, err
	}
	if space == nil {
		return nil, e.ErrNotFound
	}
	if space.UserID != userID {
		return nil, e.ErrForbidden
	}

	return taskGroup, nil
}

func (s *Service) FindAllTaskGroupsBySpace(ctx context.Context, spaceID uuid.UUID, userID uuid.UUID) ([]model.TaskGroup, error) {
	space, err := s.GetSpace(ctx, spaceID, userID)
	if err != nil {
		return nil, err
	}
	if space == nil {
		return nil, e.ErrNotFound
	}

	return s.tg.FindAllBySpace(ctx, spaceID)
}

func (s *Service) DeleteTaskGroup(ctx context.Context, id uuid.UUID) error {
	return s.tg.Delete(ctx, id)
}
