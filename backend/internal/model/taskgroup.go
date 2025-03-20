package model

import "github.com/google/uuid"

type TaskGroup struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Emoji   string    `db:"emoji" json:"emoji"`
	SpaceID uuid.UUID `db:"space_id" json:"space_id"`
}

type TaskGroupCreate struct {
	Name    string    `json:"name"`
	Emoji   string    `json:"emoji"`
	SpaceID uuid.UUID `json:"space_id"`
}
