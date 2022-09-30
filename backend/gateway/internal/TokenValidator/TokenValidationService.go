package TokenValidator

import (
	"encoding/json"
	"fmt"
	"github.com/Feokrat/music-dating-app/gateway/internal/config"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type validator struct {
	logger *log.Logger
	client http.Client
	config config.ServicesConfig
}

type ValidationService interface {
	Validate(token string) (uuid.UUID, error)
}

func NewValidationService(logger *log.Logger, config config.ServicesConfig) ValidationService {
	return validator{logger: logger, client: http.Client{}, config: config}
}

func (v validator) Validate(token string) (uuid.UUID, error) {
	sessionUrl := v.config.SessionService + fmt.Sprintf("/auth/token/validate")

	req, err := http.NewRequest("GET", sessionUrl, nil)
	if err != nil {
		v.logger.Printf("could not create request, error: %s", err.Error())
		return uuid.UUID{}, err
	}

	req.Header.Add("Authorization", "Bearer"+fmt.Sprintf(" %v", token))

	resp, err := v.client.Do(req)
	if err != nil {
		v.logger.Printf("could not get users, error: %s", err.Error())
		return uuid.UUID{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		// You may read / inspect response body
		return uuid.UUID{}, schemas.TokenError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		v.logger.Printf("could not read response body, error: %s",
			err.Error())
		return uuid.UUID{}, err
	}

	var userId schemas.IdResponse
	err = json.Unmarshal(body, &userId)
	if err != nil {
		v.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return uuid.UUID{}, err
	}

	return userId.ID, err
}
