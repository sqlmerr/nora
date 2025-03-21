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

type UserSpaceRepo interface {
	Create(ctx context.Context, userSpace *model.UserSpaceCreate) (*uuid.UUID, error)
	FindOne(ctx context.Context, id uuid.UUID) (*model.UserSpace, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserSpaceRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserSpaceRepository(dbPool *pgxpool.Pool) *UserSpaceRepository {
	return &UserSpaceRepository{
		dbPool: dbPool}
}

func (r *UserSpaceRepository) Create(ctx context.Context, userSpace *model.UserSpaceCreate) (*uuid.UUID, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("user_space").
		Columns("user_id", "space_id").
		Values(userSpace.UserID, userSpace.SpaceID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}
	var userSpaceID uuid.UUID
	result := r.dbPool.QueryRow(ctx, query, args...)

	err = result.Scan(&userSpaceID)
	return &userSpaceID, err
}

func (r *UserSpaceRepository) FindOne(ctx context.Context, id uuid.UUID) (*model.UserSpace, error) {
	query, args, err := sq.Select("*").From("user_space").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	userSpace, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.UserSpace])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userSpace, nil
}

func (r *UserSpaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete("user_space").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}
