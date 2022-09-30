package notifications

import (
	"database/sql"
	"fmt"
	"github.com/Feokrat/music-dating-app/notifications/internal/models"
	"github.com/Feokrat/music-dating-app/notifications/internal/schemas"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type messageRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

type MessageRepository interface {
	GetAllMessages(chatId uuid.UUID) ([]models.Messages, error)
	CreateMessage(message string, chatId uuid.UUID, userId uuid.UUID) (uuid.UUID, error)
}

const (
	messagesTable = "messages"
)

func NewMessageRepository(db *sqlx.DB, logger *log.Logger) MessageRepository {
	return messageRepository{
		db:     db,
		logger: logger,
	}
}

func (m messageRepository) GetAllMessages(chatId uuid.UUID) ([]models.Messages, error) {
	var messages []models.Messages
	query := fmt.Sprintf("SELECT * FROM %s WHERE chat_id = $1", messagesTable)

	err := m.db.Select(&messages, query, chatId)
	if err == sql.ErrNoRows {
		return nil, schemas.NotFoundError{Message: fmt.Sprintf("Not found any music with id %v", chatId)}
	}
	if err != nil {
		m.logger.Printf("error in db while trying to get all messages, error: %s", err.Error())
		return nil, err
	}

	return messages, nil
}

func (m messageRepository) CreateMessage(message string, chatId uuid.UUID, userId uuid.UUID) (uuid.UUID, error) {
	var messageId = uuid.New()
	query := fmt.Sprintf("INSERT INTO %s (id, creator_user_id, chat_id, content, created_at)"+
		" values ($1, $2, $3, $4, $5) RETURNING id", messagesTable)

	var id uuid.UUID

	row := m.db.QueryRow(query, messageId, userId, chatId, message, time.Now())

	if err := row.Scan(&id); err != nil {
		m.logger.Printf("error in db while trying to create message chatId:%v userId: %v, error: %s",
			chatId, userId, err.Error())
		return uuid.Nil, err
	}

	return userId, nil
}
