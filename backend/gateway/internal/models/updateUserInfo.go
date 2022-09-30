package models

type UpdateUserInfo struct {
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	HasAccess   *bool   `json:"hasAccess"`
}
