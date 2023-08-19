package handler

type Handler struct {
	logic  useCase
	logger logger.Logger
}

type useCase interface {
}

func New() Handler {
	return Handler{}
}
