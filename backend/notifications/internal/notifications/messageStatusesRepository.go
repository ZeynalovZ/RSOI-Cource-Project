package notifications

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type messageStatusesRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

type MessageStatusesRepository interface {
	TestMethod() string
}

const (
	messageStatusesTable = "messagestatuses"
)

func NewMessageStatusesRepository(db *sqlx.DB, logger *log.Logger) MessageStatusesRepository {
	return messageStatusesRepository{
		db:     db,
		logger: logger,
	}
}

func (r messageStatusesRepository) TestMethod() string {
	return "messageStatuses"
}
