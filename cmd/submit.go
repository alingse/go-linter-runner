package cmd

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/alingse/go-linter-runner/runner"
	"github.com/alingse/go-linter-runner/runner/submit"
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
		workflowRef := *workflowRefPtr
		cfg := &submit.SubmitConfig{
			Source:      sourceFile,
			RepoCount:   int(repoCount),
			Workflow:    workflow,
			WorkflowRef: workflowRef,
		}
		slog.LogAttrs(cmd.Context(), slog.LevelInfo, "submit task with", slog.Any("config", cfg))
		ctx := context.Background()
		if err := runner.Submit(ctx, cfg); err != nil {
			log.Fatal(err)
		}
		return nil
	},
}

var (
	sourceFilePtr  *string
	repoCountPtr   *int64
	workflowPtr    *string
	workflowRefPtr *string
)

func init() {
	rootCmd.AddCommand(submitCmd)
	sourceFilePtr = submitCmd.Flags().StringP("source", "s", runner.DefaultSource, "repo url file")
	repoCountPtr = submitCmd.Flags().Int64P("count", "c", runner.DefaultCount, "the repo count to submit")
	workflowPtr = submitCmd.Flags().StringP("workflow", "w", "", "workflow name to submit")
	workflowRefPtr = submitCmd.Flags().StringP("ref", "r", "", " The branch or tag name which contains the version of the workflow file you'd like to run")
}
