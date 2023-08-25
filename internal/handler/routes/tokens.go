package routes

import (
	"go-jwt-auth/internal/handler"
	"go-jwt-auth/internal/lib"
)

type TokenRoutes struct {
	tokenHandler   handler.TokenHandler
	requestHandler lib.RequestHandler
}

func NewTokenRoutes(reqHandler lib.RequestHandler, th handler.TokenHandler) TokenRoutes {
	return TokenRoutes{
		tokenHandler:   th,
		requestHandler: reqHandler,
	}
}

func (tr TokenRoutes) Setup() {
	tokens := tr.requestHandler.Gin.Group("/")
	tokens.GET("/v1/tokens", tr.tokenHandler.GetTokens)
	tokens.POST("/v1/refresh", tr.tokenHandler.RefreshTokens)
}
