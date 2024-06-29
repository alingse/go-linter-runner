package cmd

import (
	"github.com/alingse/go-linter-runner/runner"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submit repo url files into github actions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runner.RunSubmit()
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
