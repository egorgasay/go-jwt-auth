// Package repository for the possibility of using a large number of
// databases in the future through a single abstraction.
package repository

import (
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go.uber.org/fx"
)

type repository struct {
	domains.Database
	logger lib.Logger
}

var Module = fx.Options(
	fx.Provide(NewRepository),
)

func NewRepository(db domains.Database, lg lib.Logger) (domains.Repository, error) {
	return &repository{
		Database: db,
		logger:   lg,
	}, nil
}
