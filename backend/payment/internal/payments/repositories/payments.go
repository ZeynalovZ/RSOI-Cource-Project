package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Feokrat/music-dating-app/payment/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	paymentsTable = "payments"
	active        = "active"
	cancelled     = "canceled"
)

type paymentsRepository struct {
	db     *sqlx.DB
	logger *log.Logger
}

func (p paymentsRepository) CancellPayment(userId string) error {
	query := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE user_id = $2 AND status = $3`, paymentsTable)
	_, err := p.db.Exec(query, cancelled, userId, active)
	return err
}

func (p paymentsRepository) GetPaymentsByUserId(userId string) (models.Payment, error) {
	var payment []models.Payment
	var date = time.Now()
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 AND active_till_to > $2`, paymentsTable)
	err := p.db.Select(&payment, query, userId, date)
	if err == sql.ErrNoRows {
		return models.Payment{}, NotFoundError
	}
	return payment[0], err
}

func (p paymentsRepository) CreatePayment(userId string, subscriptionType int) (string, error) {
	var paymentId = uuid.New().String()
	var date = time.Now().AddDate(0, 1, 0)
	if res, err := p.CheckIfSubscriptionExists(userId, subscriptionType); err != nil {
		p.logger.Printf("error in db while trying to create payment %v, error: %s",
			paymentId, err.Error())
		return "", err
	} else {
		if res != true {
			p.logger.Printf("error in db while trying to create payment %v. Payments count not equals zero",
				paymentId)
			return "", PaymentAlreadyExists
		}
	}
	query := fmt.Sprintf(`INSERT INTO %s VALUES ($1, $2, $3, $4, $5) RETURNING id`, paymentsTable)
	row := p.db.QueryRow(query, paymentId, userId, subscriptionType, date, active)

	if err := row.Scan(&paymentId); err != nil {
		p.logger.Printf("error in db while trying to create payment %v, error: %s",
			paymentId, err.Error())
		return "", err
	}

	return paymentId, nil
}

func (p paymentsRepository) CheckIfSubscriptionExists(userId string, subscriptionType int) (bool, error) {
	var paymentsCount int
	query := fmt.Sprintf(`SELECT COUNT(*) as count FROM %s WHERE user_id = $1 AND subscription_type = $2 AND status=$3`, paymentsTable)
	err := p.db.QueryRow(query, userId, subscriptionType, active).Scan(&paymentsCount)
	if err != nil {
		p.logger.Printf("error in db while trying to check payments count %v, error: %s",
			paymentsCount, err.Error())
		return false, err
	}
	return paymentsCount == 0, err
}

type PaymentsRepository interface {
	GetPaymentsByUserId(userId string) (models.Payment, error)
	CreatePayment(userId string, subscriptionType int) (string, error)
	CancellPayment(userId string) error
}

func NewPaymentsRepository(db *sqlx.DB, logger *log.Logger) PaymentsRepository {
	return paymentsRepository{db, logger}
}
