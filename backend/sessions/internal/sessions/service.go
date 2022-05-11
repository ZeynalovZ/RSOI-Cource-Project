package sessions

import (
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/models"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/sessions/repostiroties"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/sessions/schemas"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/pkg/hash"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/pkg/token"
	"github.com/google/uuid"
	"log"
)

type AuthService struct {
	logger               *log.Logger
	credentialRepository repostiroties.CredentialsRepository
	sessionRepository    repostiroties.SessionRepository
	tokenService         token.TokenService
	hashService          hash.HashService
}

type AuthServiceInterface interface {
	SignIn(userInfo models.Auth) (string, error)
	Register(user models.Register) (string, error)
	Authorize(token string) (string, error)
}

func NewService(logger *log.Logger, credentialRepository repostiroties.CredentialsRepository,
	sessionRepository repostiroties.SessionRepository, tokenService token.TokenService,
	hashService hash.HashService) AuthService {
	return AuthService{logger: logger, credentialRepository: credentialRepository, sessionRepository: sessionRepository,
		tokenService: tokenService, hashService: hashService}
}

func (a AuthService) SignIn(userInfo models.Auth) (string, error) {
	session, err := a.sessionRepository.GetSessionByUserId(userInfo.Id)
	if err != nil {
		if err == repostiroties.NotFoundError {
			return "", schemas.InvalidCredentialsError
		}
		return "", err
	}

	credentials, err := a.credentialRepository.GetCredentialBySessionId(session.Id)
	if err != nil {
		if err == repostiroties.NotFoundError {
			return "", schemas.InvalidCredentialsError
		}
		return "", err
	}

	if err = a.hashService.ValidatePassword(userInfo.Password, credentials.PasswordHash); err != nil {
		return "", schemas.InvalidCredentialsError
	}
	return a.tokenService.GenerateToken(userInfo.Id)
}

func (a AuthService) Register(registerModel models.Register) (string, error) {
	_, err := a.credentialRepository.GetCredentialByLogin(registerModel.Login)
	if err != repostiroties.NotFoundError {
		if err == nil {
			return "", schemas.UserAlreadyExistsError
		}
		return "", err
	}
	var credential models.Credentials
	credential.Id = uuid.New().String()
	credential.SessionId = uuid.New().String()
	credential.Login = registerModel.Login
	credential.PasswordHash, err = a.hashService.HashPassword(registerModel.Password)
	if err != nil {
		return "", err
	}

	_, err = a.credentialRepository.AddCredential(credential)
	if err != nil {
		return "", err
	}

	// temp usage of user_id
	var userId = uuid.New().String()
	token, _ := a.tokenService.GenerateToken(userId)

	var session models.Sessions
	session.Id = credential.SessionId
	session.IsAuthenticated = true
	session.ExpiresAt = credential.SessionId
	session.AccessToken = token
	session.RefreshToken = token
	session.UserID = userId

	_, err = a.sessionRepository.AddSession(session)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a AuthService) Authorize(token string) (string, error) {
	userId, err := a.tokenService.ParseToken(token)
	if err != nil {
		return "", err
	}

	session, err := a.sessionRepository.GetSessionByUserId(userId)
	if err != nil {
		return "", err
	}

	if _, err = a.credentialRepository.GetCredentialBySessionId(session.Id); err != nil {
		return "", err
	}

	return userId, nil
}
