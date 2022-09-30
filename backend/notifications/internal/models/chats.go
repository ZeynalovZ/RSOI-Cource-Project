package models

import "github.com/google/uuid"

type Chats struct {
	Id      uuid.UUID `json:"id" db:"id"`
	UserId1 uuid.UUID `json:"user_id1" db:"user_id1"`
	UserId2 uuid.UUID `json:"user_id2" db:"user_id2"`
}
