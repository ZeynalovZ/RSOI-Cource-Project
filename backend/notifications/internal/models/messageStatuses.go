package models

import "github.com/google/uuid"

type MessageStatuses struct {
	Id        uuid.UUID `json:"id" db:"id"`
	MessageId uuid.UUID `json:"message_id" db:"message_id"`
	UserId    uuid.UUID `json:"user_id" db:"user_id"`
	status    bool      `json:"status" db:"status"`
}
