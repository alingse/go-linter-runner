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
		workflowName := *workflowNamePtr
		if workflowName == "" {
			return fmt.Errorf("--workflow is required")
		}
		runner.RunSubmit(sourceFile, repoCount, workflowName)
		return nil
	},
}

var sourceFilePtr *string
var repoCountPtr *int64
var workflowNamePtr *string

/*
  submit_source_file:
    description: |
      'submit' option: the repo url file for 'submit' action
      can be the file name in https://github.com/alingse/go-linter-runner/blob/main/source/ dir
      or other file link.
      Default 'top.txt'
    default: 'top.txt'
  submit_repo_count:
    description: |
      'submit' option: the repo url count to submit that read from `submit_source_file` for 'submit' action
      Default '1000'
    default: '1000'
  submit_workflow_yaml_path:
    description: |
      'submit' option: the workflow to auto
      Default '.github/workflows/run-repo.yml'
    default: '.github/workflows/run-repo.yml'

*/

func init() {
	rootCmd.AddCommand(submitCmd)
	sourceFilePtr = runCmd.Flags().StringP("source", "s", runner.DefaultSource, "repo url file")
	repoCountPtr = runCmd.Flags().Int64P("count", "c", runner.DefaultCount, "the repo count to submit")
	workflowNamePtr = runCmd.Flags().StringP("workflow", "w", runner.DefaultSource, "workflow name to submit")
}
