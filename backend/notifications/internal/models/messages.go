package models

import (
	"github.com/google/uuid"
	"time"
)

type Messages struct {
	Id            uuid.UUID `json:"id" db:"id"`
	CreatorUserId uuid.UUID `json:"creator_user_id" db:"creator_user_id"`
	ChatId        uuid.UUID `json:"chat_id" db:"chat_id"`
	Content       string    `json:"content" db:"content"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	ParentMessage uuid.UUID `json:"parent_message" db:"parent_message"`
}
