package repostiroties

import (
	"database/sql"
	"fmt"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/models"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	sessionTable = "sessions"
)

type sessionRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

func (s sessionRepository) AddSession(session models.Sessions) (string, error) {
	var sessionId string
	query := fmt.Sprintf("INSERT INTO %s (id, refresh_token, access_token, expires_at, user_id, is_authenticated)"+
		" values ($1, $2, $3, $4, $5, $6) RETURNING id", sessionTable)
	row := s.db.QueryRow(query, session.Id, session.RefreshToken, session.AccessToken, session.ExpiresAt, session.UserID, session.IsAuthenticated)
	if err := row.Scan(&sessionId); err != nil {
		return "", err
	}
	return sessionId, nil

}

func (s sessionRepository) GetSessionByUserId(userId string) (models.Sessions, error) {
	var session models.Sessions
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, sessionTable)
	err := s.db.Get(&session, query, userId)
	if err == sql.ErrNoRows {
		return session, NotFoundError
	}
	return session, err
}

type SessionRepository interface {
	GetSessionByUserId(userId string) (models.Sessions, error)
	AddSession(session models.Sessions) (string, error)
}

func NewSessionRepository(db *sqlx.DB, logger *log.Logger) SessionRepository {
	return sessionRepository{db, logger}
}
