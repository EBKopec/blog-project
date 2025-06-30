package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func respondWithError(c *gin.Context, code int, message string, err error) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
		Details: err.Error(),
	})
}
