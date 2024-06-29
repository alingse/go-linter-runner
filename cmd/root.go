package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-linter-runner",
	Short: "A cli to run linter for repos use github action",
	Long: `A cli for run linter for github repo, and it post the result as issue comment,
It can submit many check repo task by trigger many github workflow action`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
