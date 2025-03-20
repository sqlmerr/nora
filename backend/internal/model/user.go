package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"-"`
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
