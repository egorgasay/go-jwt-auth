package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"net/http"
)

type TokensHandler struct {
	tokens domains.TokenManager
	logger lib.Logger
}

func NewTokensHandler(logger lib.Logger, service domains.TokenManager) TokensHandler {
	return TokensHandler{
		logger: logger,
		tokens: service,
	}
}

func (h *TokensHandler) GetTokens(c *gin.Context) {
	guid := c.DefaultQuery("guid", "")

	access, refresh, err := h.tokens.GetTokens(c, guid)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *TokensHandler) RefreshTokens(c *gin.Context) {
	refresh := c.DefaultQuery("rtoken", "")

	access, refresh, err := h.tokens.RefreshTokens(c, refresh)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
