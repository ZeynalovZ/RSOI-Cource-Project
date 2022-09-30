package notifications

import (
	"database/sql"
	"fmt"
	"github.com/Feokrat/music-dating-app/notifications/internal/models"
	"github.com/Feokrat/music-dating-app/notifications/internal/schemas"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type chatRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

type ChatRepository interface {
	GetAllChatsByUserId(userId uuid.UUID) ([]models.Chats, error)
	CreateChat(userId1 uuid.UUID, userId2 uuid.UUID) (uuid.UUID, error)
	GetChatByUserId(userId uuid.UUID) (models.Chats, error)
}

const (
	chatTable = "chats"
)

func NewChatRepository(db *sqlx.DB, logger *log.Logger) ChatRepository {
	return chatRepository{
		db:     db,
		logger: logger,
	}
}

func (c chatRepository) GetAllChatsByUserId(userId uuid.UUID) ([]models.Chats, error) {

	var users []models.Chats
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id1 = $1 OR user_id2 = $2", chatTable)

	err := c.db.Select(&users, query, userId, userId)
	if err != nil {
		c.logger.Printf("error in db while trying to get all chats, error: %s", err.Error())
		return nil, err
	}

	return users, nil
}

func (c chatRepository) CreateChat(userId1 uuid.UUID, userId2 uuid.UUID) (uuid.UUID, error) {

	var chatId = uuid.New()
	query := fmt.Sprintf("INSERT INTO %s (id, user_id1, user_id2)"+
		" values ($1, $2, $3) RETURNING id", chatTable)

	var id uuid.UUID

	row := c.db.QueryRow(query, chatId, userId1, userId2)

	if err := row.Scan(&id); err != nil {
		c.logger.Printf("error in db while trying to create chat %v %v, error: %s",
			userId1, userId2, err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

func (c chatRepository) GetChatByUserId(userId uuid.UUID) (models.Chats, error) {
	var chat models.Chats
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, chatTable)
	err := c.db.Get(&chat, query, userId)
	if err == sql.ErrNoRows {
		return chat, schemas.NotFoundError{Message: fmt.Sprintf("Not found any chat of user with id %v", userId)}
	}

	return chat, err
}
