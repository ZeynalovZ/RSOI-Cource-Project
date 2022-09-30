package schemas

import "errors"

var (
	RightsError             = errors.New("session doesn't belong to user")
	UserAlreadyExistsError  = errors.New("user with specified login already exists")
	InvalidCredentialsError = errors.New("invalid login or password")
	TokenError              = errors.New("invalid token")
)
