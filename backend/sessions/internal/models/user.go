package models

import "github.com/google/uuid"

type Register struct {
	Login    string    `json:"login" bson:"login"`
	Password string    `json:"password" bson:"password"`
	UserId   uuid.UUID `json:"user_id"`
}

type Auth struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
}
