package models

type Sessions struct {
	Id              string `json:"id" db:"id"`
	RefreshToken    string `json:"refreshToken" db:"refresh_token"`
	AccessToken     string `json:"accessToken" db:"access_token"`
	ExpiresAt       string `json:"expiresAt" db:"expires_at"`
	IsAuthenticated bool   `json:"isAuthenticated" db:"is_authenticated"`
	UserID          string `json:"userId" db:"user_id"`
}
