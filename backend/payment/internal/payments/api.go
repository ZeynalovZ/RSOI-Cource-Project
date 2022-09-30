package payments

import (
	"github.com/Feokrat/music-dating-app/payment/internal/payments/repositories"
	"github.com/Feokrat/music-dating-app/payment/internal/payments/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	logger *log.Logger
	service PaymentsService
}

func (h handler) CreatePayment(ctx *gin.Context) {
	userIdStr := ctx.Query("user_id")
	subscriptionType := ctx.Query("subscription_type")
	_, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	i, err := strconv.Atoi(subscriptionType);
	if err != nil {
		h.logger.Printf("error has occured in subscription type conversation %v, error: %s",
			subscriptionType, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong subscription type format",
			Errors:  err.Error(),
		})

		return
	}

	paymentId, err := h.service.CreatePayment(userIdStr, i)
	if err != nil {
		switch err {
		case repositories.PaymentAlreadyExists:
			h.logger.Printf("could not create payment for user %v, error: %s",
				userIdStr, err.Error())
			ctx.JSON(http.StatusConflict, schemas.ErrorResponse{
				Message: err.Error(),
			})
		default:
			h.logger.Printf("could not create payment for user %v, error: %s",
				userIdStr, err.Error())
			ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Message: err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusCreated, schemas.PaymentResponse{PaymentId: paymentId})
}

func (h handler) GetPayment(ctx *gin.Context) {
	userIdStr := ctx.Param("user_id")
	_, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	payment, err := h.service.GetPaymentByUserId(userIdStr)
	if err != nil {
		h.logger.Printf("could not create payment for user %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, schemas.PaymentModelResponse{payment})
}

func (h handler) UpdatePayment(ctx *gin.Context) {
	userIdStr := ctx.Param("user_id")
	_, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	err = h.service.CancelPayment(userIdStr)
	if err != nil {
		h.logger.Printf("could not cancel payment for user %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

func RegisterHandlers(rg *gin.RouterGroup, service PaymentsService, logger *log.Logger) {
	h := handler{
		logger,
		service,
	}
	// sessions?test=...&m=sd
	rg.GET("/:user_id", h.GetPayment)
	rg.POST("/", h.CreatePayment)
	rg.PUT("/:user_id", h.UpdatePayment)

}