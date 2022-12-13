package models

type UpdateUserInfo struct {
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}
