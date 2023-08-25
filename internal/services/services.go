package services

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewTokenManager),
	fx.Provide(NewGeneratorService),
)
