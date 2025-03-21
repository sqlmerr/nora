package model

import "github.com/google/uuid"

type UserSpace struct {
	ID      uuid.UUID `json:"id" db:"id"`
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	SpaceID uuid.UUID `json:"space_id" db:"space_id"`
}

type UserSpaceCreate struct {
	UserID  uuid.UUID `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
}
