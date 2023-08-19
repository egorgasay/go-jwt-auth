package storage

import "go.uber.org/zap"

type Storage struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) (Storage, error) {
	return Storage{logger: logger}, nil
}

func (s Storage) Close() error {
	return nil
}
