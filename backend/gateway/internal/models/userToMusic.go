package models

import "github.com/google/uuid"

type UserToMusic struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"userId"`
	MusicId        uuid.UUID `json:"musicId"`
	FavouriteLevel int       `json:"favouriteLevel"`
}
