package routes

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/handler"
)

func Routes(r *gin.RouterGroup, h *handler.Handler) {
	r.GET("/tokens", h.GetTokens)
	r.POST("/refresh", h.RefreshTokens)
}
