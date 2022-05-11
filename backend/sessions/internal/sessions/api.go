package sessions

import (
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/models"
	"github.com/ZeynalovZ/RSOI-Course-Project/sessions/internal/sessions/schemas"
	"log"
	"net/http"
	"strings"


	"github.com/gin-gonic/gin"
)

type handler struct {
	logger  *log.Logger
	service AuthService
}

// @Summary Users sign in
// @Tags auth
// @Description Authenticate user by email and password
// @Accept json
// @Produce json
// @Param userCredentials body models.UserCredentials true "User sign in credentials"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} messageResponse
// @Failure 401 {object} messageResponse
// @Failure 500 {object} messageResponse
// @Router /auth/sign-in [post]
func signIn(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCredentials models.Auth
		if err := c.ShouldBindJSON(&userCredentials); err != nil {
			schemas.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := authService.SignIn(userCredentials)
		if err != nil {
			if err == schemas.InvalidCredentialsError {
				schemas.RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
			schemas.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		schemas.RespondWithToken(c, http.StatusOK, token)
	}
}

// @Summary Register user
// @Tags auth
// @Description Register user with specified information
// @Accept json
// @Produce json
// @Param user body models.User true "User registration information"
// @Success 201 {object} tokenResponse
// @Failure 400 {object} messageResponse
// @Failure 500 {object} messageResponse
// @Router /auth/register [post]
func register(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerModel models.Register
		if err := c.ShouldBindJSON(&registerModel); err != nil {
			schemas.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := authService.Register(registerModel)
		if err != nil {
			if err == schemas.UserAlreadyExistsError {
				schemas.RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			schemas.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		schemas.RespondWithToken(c, http.StatusCreated, token)
	}
}

func validate(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[2]

		token, err := authService.Authorize(reqToken)
		if err != nil {
			if err == schemas.UserAlreadyExistsError {
				schemas.RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			schemas.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		schemas.RespondWithToken(c, http.StatusCreated, token)
	}
}

func RegisterHandlers(rg *gin.RouterGroup, service AuthService, logger *log.Logger) {
	h := handler{
		logger:  logger,
		service: service,
	}
	rg.POST("/sign-in", signIn(h.service))
	rg.POST("/register", register(h.service))
	rg.GET("/token/validate", validate(h.service))
}
