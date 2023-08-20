package routes

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/handler"
)

func Routes(r *gin.RouterGroup, h *handler.Handler) {
	r.GET("/v1/tokens", h.GetTokens)
	r.POST("/v1/refresh", h.RefreshTokens)
}
