package bootstrap

import (
	"github.com/spf13/cobra"
	"go-jwt-auth/internal/commands"
)

var rootCmd = &cobra.Command{
	Use:              "go",
	Short:            "go-jwt-token",
	Long:             `go-jwt-token`,
	TraverseChildren: true,
}

// App root of application
type App struct {
	*cobra.Command
}

func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(commands.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
