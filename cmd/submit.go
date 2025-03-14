package cmd

import (
	"errors"
	"log"

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
			return errors.New("--source is required")
		}
		repoCount := *repoCountPtr
		if repoCount <= 0 {
			return errors.New("--count muest greater than zero")
		}
		workflow := *workflowPtr
		if workflow == "" {
			return errors.New("--workflow is required")
		}
		log.Printf("submit task with source:%s repo count: %d and workflow %s\n",
			sourceFile, repoCount, workflow)
		runner.Submit(sourceFile, repoCount, workflow)

		return nil
	},
}

var (
	sourceFilePtr *string
	repoCountPtr  *int64
	workflowPtr   *string
)

func init() {
	rootCmd.AddCommand(submitCmd)
	sourceFilePtr = submitCmd.Flags().StringP("source", "s", runner.DefaultSource, "repo url file")
	repoCountPtr = submitCmd.Flags().Int64P("count", "c", runner.DefaultCount, "the repo count to submit")
	workflowPtr = submitCmd.Flags().StringP("workflow", "w", "", "workflow name to submit")
}
