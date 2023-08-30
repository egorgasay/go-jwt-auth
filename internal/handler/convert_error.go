package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/constants"
	"net/http"
)

// HTTPError converts an error to a HTTP error
func HTTPError(c *gin.Context, err error) {
	switch err { // no errors.Is() because we get an explicit error from the service every time.
	case constants.ErrMissingRefreshToken, constants.ErrMissingAccessToken:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	case constants.ErrTokenExpired:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	case constants.ErrInvalidToken, constants.ErrInvalidGUID, constants.ErrNotFound:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
