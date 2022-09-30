package models

import "github.com/google/uuid"

type Image struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"userId"`
	Image  string    `json:"image"`
}
