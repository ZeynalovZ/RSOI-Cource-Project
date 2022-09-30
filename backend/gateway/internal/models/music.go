package models

import "github.com/google/uuid"

type Music struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Author string    `json:"author"`
	Url    string    `json:"url"`
}
