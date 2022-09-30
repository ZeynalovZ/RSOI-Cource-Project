package schemas

import (
	"github.com/Feokrat/music-dating-app/notifications/internal/models"
	"github.com/google/uuid"
)

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

type MessageRequest struct {
	UserId  uuid.UUID `json:"user_id"`
	ChatId  uuid.UUID `json:"chat_id"`
	Message string    `json:"message"`
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type ChatsResponse struct {
	Chats []ChatsModel
}

type ChatsModel struct {
	Id          uuid.UUID `json:"id"`
	LastMessage string    `json:"lastMessage"`
	IsRead      bool      `json:"isRead"`
	UserId1     uuid.UUID `json:"UserId1"`
	UserId2     uuid.UUID `json:"UserId2"`
}

type MessageResponse struct {
	Messages []models.Messages `json:"messages"`
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
