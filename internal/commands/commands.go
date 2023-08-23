package commands

import (
	"context"
	"github.com/spf13/cobra"
	"go-jwt-auth/internal/lib"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log"
)

var cmds = map[string]lib.Command{
	"go": NewGoCommand(),
}

func WrapSubCommand(name string, cmd lib.Command, opt fx.Option) *cobra.Command {
	wrappedCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			logger, err := lib.NewSugaredLogger()
			if err != nil {
				log.Fatalln(err)
			}

			opts := fx.Options(
				fx.WithLogger(func() fxevent.Logger {
					return logger
				}),
				fx.Invoke(cmd.Run()),
			)
			ctx := context.Background()
			app := fx.New(opt, opts)
			err = app.Start(ctx)
			defer app.Stop(ctx)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}
	cmd.Setup(wrappedCmd)
	return wrappedCmd
}
