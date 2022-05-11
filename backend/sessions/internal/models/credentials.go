package models

type Credentials struct {
	Id           string `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"passwordHash" db:"password_hash"`
	SessionId    string `json:"sessionId" db:"session_id"`
	Role         string `json:"role" db:"role"`
}