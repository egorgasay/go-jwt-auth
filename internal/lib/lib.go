package lib

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewRequestHandler),
	fx.Provide(NewDatabase),
	fx.Provide(NewLogger),
)
