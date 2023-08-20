package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logic  useCase
	logger *zap.Logger
}

type useCase interface {
	GetTokens(guid string) (string, string, error)
	RefreshTokens(refresh string) (string, string, error)
}

func New(logic useCase, logger *zap.Logger) *Handler {
	return &Handler{logic: logic, logger: logger}
}

func (h *Handler) GetTokens(c *gin.Context) {
	guid := c.DefaultQuery("guid", "")

	access, refresh, err := h.logic.GetTokens(guid)
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
	refresh := c.DefaultQuery("refresh", "")

	access, refresh, err := h.logic.RefreshTokens(refresh)
	if err != nil {
		HTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
