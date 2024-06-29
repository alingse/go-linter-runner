package cmd

import (
	"fmt"
	"os"

	"github.com/alingse/go-linter-runner/runner"
	"github.com/alingse/go-linter-runner/runner/utils"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check the repo by a given config",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonConfig := readFileOrString(*jsonConfigPtr)
		yamlConfig := readFileOrString(*yamlConfigPtr)
		if jsonConfig == "" && yamlConfig == "" {
			return fmt.Errorf("one of the options -y and -j must be set.")
		}
		repo := *repoURLPtr
		if repo == "" {
			return fmt.Errorf("the -r/--repo muest be set.")
		}
		runner.RunCheck()
		return nil
	},
}

var jsonConfigPtr *string
var yamlConfigPtr *string
var repoURLPtr *string

func init() {
	rootCmd.AddCommand(checkCmd)
	jsonConfigPtr = checkCmd.Flags().StringP("-json", "j", "", "a json config file or a json config string")
	yamlConfigPtr = checkCmd.Flags().StringP("-yaml", "y", "", "a yaml config file or a yaml config string")
	repoURLPtr = checkCmd.Flags().StringP("repo", "r", "", "the repo needs to lint")
}

func readFileOrString(f string) string {
	if f == "" {
		return ""
	}
	if utils.IsFileExists(f) {
		data, err := os.ReadFile(f)
		if err != nil {
			return ""
		}
		return string(data)
	}
	return f
}
