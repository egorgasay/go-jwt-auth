package handler

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"net/http"
	"strings"
)

type TokenHandler struct {
	tokens domains.TokenManager
	logger lib.Logger
}

func NewTokenHandler(logger lib.Logger, service domains.TokenManager) TokenHandler {
	return TokenHandler{
		logger: logger,
		tokens: service,
	}
}

func (h *TokenHandler) GetTokens(c *gin.Context) {
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

func (h *TokenHandler) RefreshTokens(c *gin.Context) {
	access := c.Request.Header.Get("Authorization")
	access = strings.TrimLeft(access, "Bearer ")

	rtr := &RefreshTokensRequest{}
	if err := c.BindJSON(rtr); err != nil {
		HTTPError(c, err)
		return
	}

	access, refresh, err := h.tokens.RefreshTokens(c, access, rtr.RefreshToken)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
