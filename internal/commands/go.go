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
		middleware middlewares.Middlewares,
		env lib.Env,
		router lib.RequestHandler,
		route routes.Routes,
		logger lib.Logger,
		database lib.Database,
	) {
		middleware.Setup()
		route.Setup()

		logger.Info("Running server")
		if env.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + env.ServerPort)
		}
	}
}

func NewGoCommand() *GoCommand {
	return &GoCommand{}
}
