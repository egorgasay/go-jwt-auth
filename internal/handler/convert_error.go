package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/usecase"
	"net/http"
)

func HTTPError(c *gin.Context, err error) {
	switch err { // no errors.Is() because we get an explicit error from the usecase every time.
	case usecase.ErrNoToken:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	case usecase.ErrInvalidToken:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	case usecase.ErrSign, usecase.ErrGenerateUUID:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
