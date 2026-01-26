package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends a JSON error response
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{Error: message})
}

// HandleDBError handles common database errors with appropriate HTTP status codes
func HandleDBError(c *gin.Context, err error, resourceName string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		RespondWithError(c, http.StatusNotFound, resourceName+" not found")
		return
	}
	RespondWithError(c, http.StatusInternalServerError, err.Error())
}
