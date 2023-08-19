package handler

import "go.uber.org/zap"

type Handler struct {
	logic  useCase
	logger *zap.Logger
}

type useCase interface {
}

func New(logic useCase, logger *zap.Logger) Handler {
	return Handler{logic: logic, logger: logger}
}
