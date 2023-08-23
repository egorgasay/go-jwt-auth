package lib

import "github.com/spf13/cobra"

type CommandRunner any

type Command interface {
	Short() string
	Setup(cmd *cobra.Command)
	Run() CommandRunner
}
