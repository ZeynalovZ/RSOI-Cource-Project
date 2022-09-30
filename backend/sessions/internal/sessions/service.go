package sessions

import (
	"github.com/Feokrat/music-dating-app/sessions/internal/models"
	"github.com/Feokrat/music-dating-app/sessions/internal/sessions/schemas"
	"log"

	"github.com/Feokrat/music-dating-app/sessions/internal/sessions/repostiroties"

	"github.com/Feokrat/music-dating-app/sessions/pkg/hash"
	"github.com/Feokrat/music-dating-app/sessions/pkg/token"
	"github.com/google/uuid"
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

	credentials, err := a.credentialRepository.GetCredentialByLogin(userInfo.Login)
	if err != nil {
		if err == repostiroties.NotFoundError {
			return "", schemas.InvalidCredentialsError
		}
		return "", err
	}

	if err = a.hashService.ValidatePassword(userInfo.Password, credentials.PasswordHash); err != nil {
		return "", schemas.InvalidCredentialsError
	}

	session, err := a.sessionRepository.GetSessionById(credentials.SessionId)
	if err == repostiroties.NotFoundError {
		return "", schemas.InvalidCredentialsError
	}

	return session.AccessToken, nil
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
	var userId = registerModel.UserId
	token, _ := a.tokenService.GenerateToken(userId.String())

	var session models.Sessions
	session.Id = credential.SessionId
	session.IsAuthenticated = true
	// TODO: real expiration here is needed
	session.ExpiresAt = "10-12-2023"
	session.AccessToken = token
	session.RefreshToken = token
	session.UserID = userId.String()

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
