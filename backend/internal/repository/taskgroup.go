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

type TaskGroupRepo interface {
	Create(ctx context.Context, taskGroup *model.TaskGroupCreate) (*uuid.UUID, error)
	FindOne(ctx context.Context, id uuid.UUID) (*model.TaskGroup, error)
	FindAllBySpace(ctx context.Context, spaceID uuid.UUID) ([]model.TaskGroup, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskGroupRepository struct {
	dbPool *pgxpool.Pool
}

func NewTaskGroupRepository(dbPool *pgxpool.Pool) *TaskGroupRepository {
	return &TaskGroupRepository{
		dbPool: dbPool}
}

func (r *TaskGroupRepository) Create(ctx context.Context, taskGroup *model.TaskGroupCreate) (*uuid.UUID, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("taskgroups").
		Columns("name", "emoji", "space_id").
		Values(taskGroup.Name, taskGroup.Emoji, taskGroup.SpaceID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}
	var taskGroupID uuid.UUID
	result := r.dbPool.QueryRow(ctx, query, args...)

	err = result.Scan(&taskGroupID)
	return &taskGroupID, err
}

func (r *TaskGroupRepository) FindOne(ctx context.Context, id uuid.UUID) (*model.TaskGroup, error) {
	query, args, err := sq.Select("*").From("taskgroups").Where("id = $1", id).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	taskGroup, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.TaskGroup])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &taskGroup, nil
}

func (r *TaskGroupRepository) FindAllBySpace(ctx context.Context, spaceID uuid.UUID) ([]model.TaskGroup, error) {
	query, args, err := sq.Select("*").From("taskgroups").Where("space_id = $1", spaceID).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	taskGroups, err := pgx.CollectRows(result, pgx.RowToStructByName[model.TaskGroup])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return taskGroups, nil
}

func (r *TaskGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete("taskgroups").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}
