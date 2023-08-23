package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/services"
	"net/http"
)

func HTTPError(c *gin.Context, err error) {
	switch err { // no errors.Is() because we get an explicit error from the service every time.
	case services.ErrNoToken:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	case services.ErrExpired:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	case services.ErrInvalidToken, services.ErrInvalidGUID:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	case services.ErrSign, services.ErrGenerateUUID, services.ErrCantHashToken:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
