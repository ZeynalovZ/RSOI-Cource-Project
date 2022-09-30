package models

import "github.com/google/uuid"

type Image struct {
	Id     uuid.UUID `json:"id" db:"id"`
	UserId uuid.UUID `json:"userId" db:"user_id"`
	Image  string    `json:"image" db:"image"`
}
