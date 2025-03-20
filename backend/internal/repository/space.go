package repository

import (
	"context"
	"errors"
	"nora/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpaceRepo interface {
	Create(ctx context.Context, space *model.SpaceCreate) (*uuid.UUID, error)
	FindOne(ctx context.Context, id uuid.UUID) (*model.Space, error)
	FindAllByUser(ctx context.Context, userID uuid.UUID) ([]model.Space, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SpaceRepository struct {
	dbPool *pgxpool.Pool
}

func NewSpaceRepository(dbPool *pgxpool.Pool) *SpaceRepository {
	return &SpaceRepository{
		dbPool: dbPool}
}

func (r *SpaceRepository) Create(ctx context.Context, space *model.SpaceCreate) (*uuid.UUID, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("spaces").
		Columns("name", "user_id").
		Values(space.Name, space.UserID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}
	var spaceID uuid.UUID
	result := r.dbPool.QueryRow(ctx, query, args...)

	err = result.Scan(&spaceID)
	return &spaceID, err
}

func (r *SpaceRepository) FindOne(ctx context.Context, id uuid.UUID) (*model.Space, error) {
	query, args, err := sq.Select("*").From("spaces").Where("id = $1", id).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	space, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.Space])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &space, nil
}

func (r *SpaceRepository) FindAllByUser(ctx context.Context, userID uuid.UUID) ([]model.Space, error) {
	query, args, err := sq.Select("*").From("spaces").Where("user_id = $1", userID).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	spaces, err := pgx.CollectRows(result, pgx.RowToStructByName[model.Space])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return spaces, nil
}

func (r *SpaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete("spaces").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}
