package cmd

import (
	"errors"
	"os"

	"github.com/alingse/go-linter-runner/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the linter for repo by a given config",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		yamlData, _ := os.ReadFile(*yamlConfigPtr)
		yamlConfig := string(yamlData)
		jsonConfig := *jsonConfigPtr
		if yamlConfig == "" && jsonConfig == "" {
			return errors.New("one of the options -y and -j must be set.")
		}

		repo := *repoURLPtr
		if repo == "" {
			return errors.New("the -r/--repo muest be set.")
		}
		runner.Run(repo, jsonConfig, yamlConfig)

		return nil
	},
}

var (
	jsonConfigPtr *string
	yamlConfigPtr *string
	repoURLPtr    *string
)

func init() {
	rootCmd.AddCommand(runCmd)
	yamlConfigPtr = runCmd.Flags().StringP("yaml", "y", "", "a yaml config file")
	jsonConfigPtr = runCmd.Flags().StringP("json", "j", "", "a json string config")
	repoURLPtr = runCmd.Flags().StringP("repo", "r", "", "the repo needs to lint")
}
