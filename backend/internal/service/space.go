package service

import (
	"context"
	e "nora/internal/error"
	"nora/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) CreateSpace(ctx context.Context, data *model.SpaceCreate) (*model.Space, error) {
	id, err := s.sp.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	logger := ctx.Value("logger").(*zap.Logger)
	logger.Info("created space", zap.String("id", id.String()))

	return s.sp.FindOne(ctx, *id)
}

func (s *Service) GetSpace(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Space, error) {
	space, err := s.sp.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if space == nil {
		return nil, e.ErrNotFound
	}
	if space.UserID != userID {
		return nil, e.ErrForbidden
	}
	return space, nil
}

func (s *Service) ListSpaces(ctx context.Context, userID uuid.UUID) ([]model.Space, error) {
	return s.sp.FindAllByUser(ctx, userID)
}

func (s *Service) DeleteSpace(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	space, err := s.sp.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if space == nil {
		return e.ErrNotFound
	}
	if space.UserID != userID {
		return e.ErrForbidden
	}
	logger := ctx.Value("logger").(*zap.Logger)
	logger.Info("space deleted", zap.String("id", id.String()))
	return s.sp.Delete(ctx, id)
}
