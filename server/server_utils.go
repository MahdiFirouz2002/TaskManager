package server

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func RespondWithError(c *gin.Context, statusCode int, message string, code string) {
	c.JSON(statusCode, ErrorResponse{
		Error: message,
		Code:  code,
	})
}
