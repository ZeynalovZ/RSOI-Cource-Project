package models

import "github.com/google/uuid"

type UserToMusic struct {
	Id             uuid.UUID `json:"id" db:"id"`
	UserId         uuid.UUID `json:"userId" db:"user_id"`
	MusicId        uuid.UUID `json:"musicId" db:"music_id"`
	FavouriteLevel int       `json:"favouriteLevel" db:"favourite_level"`
}
