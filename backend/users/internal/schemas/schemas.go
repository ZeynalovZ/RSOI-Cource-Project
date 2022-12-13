package schemas

import (
	"github.com/Feokrat/music-dating-app/users/internal/models"
	"github.com/google/uuid"
)

type UserImageResponse struct {
	Image string `json:"image"`
}

type LikeRequest struct {
	UserId uuid.UUID `json:"userId"`
}

type UserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	HasAccess   bool   `json:"hasAccess"`
}

type UpdateRequest struct {
	Name    *string `json:"name" db:"name"`
	Surname *string `json:"surname" db:"surname"`
	Image   string  `json:"image" db:"image"`
}

type MusicRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Url    string `json:"url"`
}

type UserToMusicRequest struct {
	UserId         uuid.UUID `json:"userId" db:"user_id"`
	MusicId        uuid.UUID `json:"musicId" db:"music_id"`
	FavouriteLevel int       `json:"favouriteLevel" db:"favourite_level"`
}

type ImageRequest struct {
	UserId uuid.UUID `json:"userId" db:"user_id"`
	Image  string    `json:"image" db:"image"`
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type LikeResponse struct {
	IsMatch bool
}

type UserResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Description      string    `json:"description"`
	Image            string    `json:"image"`
	SubscriptionType int       `json:"subscriptionType"`
	MusicIds         []string  `json:"musicIds"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

type MusicResponse struct {
	Music models.Music `json:"music"`
}

type MusicsResponse struct {
	Musics []models.Music `json:"musics"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string `json:"message"`
}

func (e NotFoundError) Error() string {
	return e.Message
}
