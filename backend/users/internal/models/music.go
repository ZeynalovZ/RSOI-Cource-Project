package models

import "github.com/google/uuid"

type Music struct {
	Id     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Author string    `json:"author" db:"author"`
	Url    string    `json:"url" db:"url"`
}
