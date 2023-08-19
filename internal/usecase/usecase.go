package usecase

import "go.uber.org/zap"

type UseCase struct {
	storage storage
	logger  *zap.Logger
}

type storage interface {
}

func New(st storage, logger *zap.Logger) UseCase {
	return UseCase{storage: st, logger: logger}
}
