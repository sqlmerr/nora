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

type UserRepo interface {
	Create(ctx context.Context, user *model.UserCreate) (*uuid.UUID, error)
	FindOne(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindOneByUsername(ctx context.Context, username string) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, data *model.UserUpdate) error
}

type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		dbPool: dbPool}
}

func (r *UserRepository) Create(ctx context.Context, user *model.UserCreate) (*uuid.UUID, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("users").
		Columns("name", "username", "password").
		Values(user.Name, user.Username, user.Password).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}
	var userID uuid.UUID
	result := r.dbPool.QueryRow(ctx, query, args...)

	result.Scan(&userID)

	return &userID, nil
}

func (r *UserRepository) FindOne(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query, args, err := sq.Select("*").From("users").Where("id = $1", id).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindOneByUsername(ctx context.Context, username string) (*model.User, error) {
	query, args, err := sq.Select("*").From("users").Where("username = $1", username).ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(result, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete("users").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, data *model.UserUpdate) error {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update("users")
	if data.Name != "" {
		b = b.Set("name", data.Name)
	}
	if data.Username != "" {
		b = b.Set("username", data.Username)
	}
	if data.Password != "" {
		b = b.Set("password", data.Password)
	}
	if data.TelegramID != 0 {
		b = b.Set("telegram_id", data.TelegramID)
	}

	query, args, err := b.Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}
