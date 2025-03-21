package model

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"-"`
	TelegramID int64     `db:"telegram_id" json:"telegram_id"`
}

type UserCreate struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	TelegramID int64  `json:"telegram_id"`
}

type UserTelegramConnect struct {
	TelegramID int64     `json:"telegram_id"`
	UserID     uuid.UUID `json:"user_id"`
}
