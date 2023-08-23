package storage

import (
	"go.uber.org/fx"
)

type Config struct {
	DatabaseDSN string `json:"dsn"`
}

var Module = fx.Options(
	fx.Provide(NewDatabase),
)
