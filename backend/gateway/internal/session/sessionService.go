package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Feokrat/music-dating-app/gateway/internal/config"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/Feokrat/music-dating-app/gateway/internal/session/models"
	"io"
	"log"
	"net/http"
)

type service struct {
	logger *log.Logger
	client http.Client
	config config.ServicesConfig
}

func (s service) Authorize(auth models.Auth) (schemas.TokenResponse, error) {
	getAuthUrl := s.config.SessionService + fmt.Sprintf("/auth/sign-in")

	var authBytes bytes.Buffer
	err := json.NewEncoder(&authBytes).Encode(auth)
	if err != nil {
		s.logger.Printf("could not convert to io read auth, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	req, err := http.NewRequest("POST", getAuthUrl, &authBytes)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get users, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.TokenResponse{}, err
	}

	var token schemas.TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	return token, err
}

func (s service) Register(auth models.Register) (schemas.TokenResponse, error) {
	getAuthUrl := s.config.SessionService + fmt.Sprintf("/auth/register")

	var authBytes bytes.Buffer
	err := json.NewEncoder(&authBytes).Encode(auth)
	if err != nil {
		s.logger.Printf("could not convert to io read auth, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	req, err := http.NewRequest("POST", getAuthUrl, &authBytes)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get users, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		// You may read / inspect response body
		return schemas.TokenResponse{}, schemas.UserAlreadyExistsError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.TokenResponse{}, err
	}

	var token schemas.TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.TokenResponse{}, err
	}

	return token, err
}

type SessionService interface {
	Authorize(auth models.Auth) (schemas.TokenResponse, error)
	Register(auth models.Register) (schemas.TokenResponse, error)
}

func NewSessionService(logger *log.Logger, config config.ServicesConfig) SessionService {
	return service{logger: logger, client: http.Client{}, config: config}
}
