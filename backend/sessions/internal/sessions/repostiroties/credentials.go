package repostiroties

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Feokrat/music-dating-app/sessions/internal/models"
	"github.com/jmoiron/sqlx"
)

type credentialsRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

const (
	credentialsTable = "credentials"
	adminRole        = "admin"
	userRole         = "user"
	primaryUserRole  = "primary_user"
)

func (c credentialsRepository) AddCredential(credential models.Credentials) (string, error) {
	var credentialId string
	query := fmt.Sprintf("INSERT INTO %s (id, login, password_hash, session_id, role)"+
		" values ($1, $2, $3, $4, $5) RETURNING id", credentialsTable)
	row := c.db.QueryRow(query, credential.Id, credential.Login, credential.PasswordHash, credential.SessionId, userRole)
	if err := row.Scan(&credentialId); err != nil {
		return "", err
	}
	return credentialId, nil
}

func (c credentialsRepository) GetCredentialBySessionId(sessionId string) (models.Credentials, error) {
	var credential models.Credentials
	query := fmt.Sprintf(`SELECT * FROM %s WHERE session_id = $1`, credentialsTable)
	err := c.db.Get(&credential, query, sessionId)
	if err == sql.ErrNoRows {
		return credential, NotFoundError
	}
	return credential, err
}

func (c credentialsRepository) GetCredentialByLogin(login string) (models.Credentials, error) {
	var credential models.Credentials
	query := fmt.Sprintf(`SELECT * FROM %s WHERE login = $1`, credentialsTable)
	err := c.db.Get(&credential, query, login)
	if err == sql.ErrNoRows {
		return credential, NotFoundError
	}
	return credential, err
}

type CredentialsRepository interface {
	GetCredentialBySessionId(sessionId string) (models.Credentials, error)
	GetCredentialByLogin(login string) (models.Credentials, error)
	AddCredential(credential models.Credentials) (string, error)
}

func NewCredentialsRepository(db *sqlx.DB, logger *log.Logger) CredentialsRepository {
	return credentialsRepository{db, logger}
}
