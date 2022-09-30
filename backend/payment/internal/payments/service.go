package payments

import (
	"github.com/Feokrat/music-dating-app/payment/internal/models"
	"github.com/Feokrat/music-dating-app/payment/internal/payments/repositories"
	"log"
)

type paymentsService struct {
	logger *log.Logger
	paymentsRepository repositories.PaymentsRepository
}

func (p paymentsService) CreatePayment(userId string, subscriptionType int) (string, error) {
	// получить все
	// проверить, если есть закенселенная, но работающая, то изменить статус
	// иначе создать новую
	return p.paymentsRepository.CreatePayment(userId, subscriptionType)
}

func (p paymentsService) CancelPayment(userId string) error {
	return p.paymentsRepository.CancellPayment(userId)
}

func (p paymentsService) GetPaymentByUserId(userId string) (models.Payment, error) {
	return p.paymentsRepository.GetPaymentsByUserId(userId)
}

type PaymentsService interface {
	CreatePayment(userId string, subscriptionType int) (string, error)
	CancelPayment(userId string) error
	GetPaymentByUserId(userId string) (models.Payment, error)
}

func NewPaymentService(logger *log.Logger, paymentsRepository repositories.PaymentsRepository) PaymentsService {
	return paymentsService{logger, paymentsRepository}
}