package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/service"
	"net/http"
)

func HTTPError(c *gin.Context, err error) {
	switch err { // no errors.Is() because we get an explicit error from the service every time.
	case service.ErrNoToken, service.ErrExpired:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	case service.ErrInvalidToken, service.ErrInvalidGUID:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	case service.ErrSign, service.ErrGenerateUUID:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
