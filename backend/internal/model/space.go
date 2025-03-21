package model

import "github.com/google/uuid"

type Space struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Name   string    `db:"name" json:"name"`
	UserID uuid.UUID `db:"user_id" json:"user_id"` // owner id
}

type SpaceCreate struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"-"` // owner id
}
