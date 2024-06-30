package cmd

import (
	"fmt"

	"github.com/alingse/go-linter-runner/runner"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submit repo url files into github actions",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceFile := *sourceFilePtr
		if sourceFile == "" {
			return fmt.Errorf("--source is required")
		}
		repoCount := *repoCountPtr
		if repoCount <= 0 {
			return fmt.Errorf("--count muest greater than zero")
		}
		workflow := *workflowPtr
		if workflow == "" {
			return fmt.Errorf("--workflow is required")
		}
		runner.RunSubmit(sourceFile, repoCount, workflow)
		return nil
	},
}

var sourceFilePtr *string
var repoCountPtr *int64
var workflowPtr *string

func init() {
	rootCmd.AddCommand(submitCmd)
	sourceFilePtr = runCmd.Flags().StringP("source", "s", runner.DefaultSource, "repo url file")
	repoCountPtr = runCmd.Flags().Int64P("count", "c", runner.DefaultCount, "the repo count to submit")
	workflowPtr = runCmd.Flags().StringP("workflow", "w", runner.DefaultSource, "workflow name to submit")
}
