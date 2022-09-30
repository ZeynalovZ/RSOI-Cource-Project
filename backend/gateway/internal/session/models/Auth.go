package models

import "github.com/google/uuid"

type Auth struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
}

type Register struct {
	UserId   uuid.UUID `json:"user_id"`
	Login    string    `json:"login" bson:"login"`
	Password string    `json:"password" bson:"password"`
}
