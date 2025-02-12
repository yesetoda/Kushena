package helper

import (
	"github.com/gin-gonic/gin"

	"github.com/yesetoda/Kushena/models"
)

// Success Response
func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, models.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Error Response
func ErrorResponse(c *gin.Context, status int, errMsg string) {
	c.JSON(status, models.Response{
		Status:  status,
		Message: "Request failed",
		Error:   errMsg,
	})
}
