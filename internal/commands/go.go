package commands

import (
	"github.com/spf13/cobra"
	"go-jwt-auth/internal/handler/routes"
	"go-jwt-auth/internal/lib"
)

type GoCommand struct{}

func (s *GoCommand) Short() string {
	return "serve application"
}

func (s *GoCommand) Setup(cmd *cobra.Command) {}

func (s *GoCommand) Run() lib.CommandRunner {
	return func(
		conf lib.Config,
		route routes.Routes,
		reqHandler lib.RequestHandler,
		logger lib.Logger,
		database lib.Database,
	) {
		route.Setup()

		logger.Info("Running server")
		err := reqHandler.Gin.Run(":" + conf.Port)
		if err != nil {
			return
		}
	}
}

func NewGoCommand() *GoCommand {
	return &GoCommand{}
}
