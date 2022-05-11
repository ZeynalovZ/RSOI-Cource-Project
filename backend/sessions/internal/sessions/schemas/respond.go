package schemas

import "github.com/gin-gonic/gin"

type dataResponse struct {
	Data interface{} `json:"data"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, messageResponse{message})
}

func RespondWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, dataResponse{data})
}

func RespondWithId(c *gin.Context, statusCode int, id int) {
	c.JSON(statusCode, idResponse{id})
}

func RespondWithToken(c *gin.Context, statusCode int, token string) {
	c.JSON(statusCode, idResponse{token})
}
