package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logic  useCase
	logger *zap.Logger
}

type useCase interface {
	GetTokens(ctx context.Context, guid string) (string, string, error)
	RefreshTokens(ctx context.Context, refresh string) (string, string, error)
}

func New(logic useCase, logger *zap.Logger) *Handler {
	return &Handler{logic: logic, logger: logger}
}

func (h *Handler) GetTokens(c *gin.Context) {
	guid := c.DefaultQuery("guid", "")

	access, refresh, err := h.logic.GetTokens(c, guid)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	refresh := c.DefaultQuery("rtoken", "")

	access, refresh, err := h.logic.RefreshTokens(c, refresh)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
