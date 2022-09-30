package session

import (
	"github.com/Feokrat/music-dating-app/gateway/internal/gateway"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/Feokrat/music-dating-app/gateway/internal/session/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type handler struct {
	logger      *log.Logger
	service     SessionService
	userService gateway.UsersService
}

func (h handler) Authorize(ctx *gin.Context) {
	var userCredentials models.Auth
	if err := ctx.ShouldBindJSON(&userCredentials); err != nil {
		schemas.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	answ, err := h.service.Authorize(userCredentials)
	if err != nil {
		h.logger.Println("authorization failed for %v", userCredentials.Login)
		ctx.JSON(http.StatusUnauthorized, schemas.Error500response{Message: "authorization failed for user", Code: 500})
	}

	ctx.JSON(http.StatusOK, answ)
}

func (h handler) Register(ctx *gin.Context) {
	var userCredentials models.Register
	if err := ctx.ShouldBindJSON(&userCredentials); err != nil {
		schemas.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var user = schemas.UserRequest{}
	id, err := h.userService.AddUser(user)
	if err != nil {
		h.logger.Printf("could not add user %v, error: %s",
			user, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	userCredentials.UserId = id

	answ, err := h.service.Register(userCredentials)
	if err != nil {
		h.logger.Println("registration failed for %v", userCredentials.Login)
		ctx.JSON(http.StatusUnauthorized, schemas.Error500response{Message: "registration failed for user", Code: 500})
	}

	ctx.JSON(http.StatusOK, answ)
}

func RegisterAuthHandlers(rg *gin.RouterGroup, service SessionService, logger *log.Logger, userService gateway.UsersService) {
	h := handler{logger: logger, service: service, userService: userService}
	rg.POST("/login", h.Authorize)
	rg.POST("/register", h.Register)
}
