package models

type UpdateUserInfo struct {
	Name        *string `json:"name" db:"name"`
	Surname     *string `json:"surname" db:"surname"`
	Email       *string `json:"email" db:"email"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	HasAccess   *bool   `json:"hasAccess" db:"has_access"`
}
