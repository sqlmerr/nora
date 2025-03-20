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

type TaskRepo interface {
	Create(ctx context.Context, task *model.TaskCreate) (*uuid.UUID, error)
	FindOne(ctx context.Context, id uuid.UUID) (*model.Task, error)
	FindAllByGroup(ctx context.Context, groupID uuid.UUID) ([]model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, task *model.TaskUpdate) error
}

type TaskRepository struct {
	dbPool *pgxpool.Pool
}

func NewTaskRepository(dbPool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		dbPool: dbPool}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.TaskCreate) (*uuid.UUID, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("tasks").
		Columns("name", "description", "group_id").
		Values(task.Name, task.Description, task.GroupID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}
	var taskID uuid.UUID
	result := r.dbPool.QueryRow(ctx, query, args...)

	err = result.Scan(&taskID)
	return &taskID, err
}

func (r *TaskRepository) FindOne(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	query, args, err := sq.Select("*").From("tasks").Where("id = $1", id).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	task, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) FindAllByGroup(ctx context.Context, groupID uuid.UUID) ([]model.Task, error) {
	query, args, err := sq.Select("*").From("tasks").Where("group_id = $1", groupID).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	tasks, err := pgx.CollectRows(result, pgx.RowToStructByName[model.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete("tasks").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, task *model.TaskUpdate) error {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update("tasks")
	if task.Name != "" {
		b = b.Set("name", task.Name)
	}
	if task.Description != "" {
		b = b.Set("description", task.Description)
	}
	if task.GroupID != uuid.Nil {
		b = b.Set("group_id", task.GroupID)
	}

	query, args, err := b.Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}
