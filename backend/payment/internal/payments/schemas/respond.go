package schemas

import "github.com/Feokrat/music-dating-app/payment/internal/models"

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type PaymentResponse struct {
	PaymentId string `json:"paymentId"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string `json:"message"`
}

type PaymentModelResponse struct {
	Payment models.Payment `json:"payment"`
}