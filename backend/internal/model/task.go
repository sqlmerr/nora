package model

import "github.com/google/uuid"

type Task struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	GroupID     uuid.UUID `db:"group_id" json:"group_id"`
}

type TaskCreate struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"group_id"`
}

type TaskUpdate struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"group_id"`
}
