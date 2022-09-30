package models

import "github.com/google/uuid"

type UserLikes struct {
	Id      uuid.UUID `json:"id" db:"id"`
	Who     uuid.UUID `json:"who" db:"who"`
	FromWho uuid.UUID `json:"fromWho" db:"from_who"`
}
