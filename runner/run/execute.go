package run

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/alingse/go-linter-runner/runner/utils"
)

var ErrSkipNoGoModRepo = errors.New("skip this repo for no go.mod file exists")

func runCmd(cmd *exec.Cmd) error {
	data, err := cmd.CombinedOutput()
	fmt.Println(string(data))
	if err != nil {
		return fmt.Errorf("run %s %+v failed %w", cmd.Path, cmd.Args, err)
	}
	return nil
}

func Prepare(ctx context.Context, cfg *Config) error {
	// install linter
	name, args := utils.SplitCommand(cfg.LinterCfg.InstallCommand)
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = cfg.LinterCfg.Workdir
	if err := runCmd(cmd); err != nil {
		return err
	}

	// clone repo
	cmd = exec.CommandContext(ctx, "git", "clone", cfg.Repo)
	cmd.Dir = cfg.LinterCfg.Workdir
	if err := runCmd(cmd); err != nil {
		return err
	}

	// TODO: check more deep
	// check go.mod exists
	gomodFile := path.Join(cfg.RepoDir, "go.mod")
	if !utils.IsFileExists(gomodFile) {
		return ErrSkipNoGoModRepo
	}

	// run go mod download
	cmd = exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = cfg.RepoDir
	if err := runCmd(cmd); err != nil {
		return err
	}

	// read default branch for repo
	cmd = exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = cfg.RepoDir
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git branch failed %w", err)
	}
	cfg.RepoBranch = strings.TrimSpace(string(output))
	return nil
}

func Run(ctx context.Context, cfg *Config) ([]string, error) {
	name, args := utils.SplitCommand(cfg.LinterCfg.LinterCommand)
	args = append(args, "./...")
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = cfg.RepoDir
	data, err := cmd.CombinedOutput()
	output := string(data)
	if err != nil && len(output) == 0 {
		return nil, err
	}

	// check includes && excludes
	outputs := strings.Split(output, "\n")
	validOutputs := make([]string, 0, len(outputs))

	includes := utils.GetStringArray(cfg.LinterCfg.Includes)
	excludes := utils.GetStringArray(cfg.LinterCfg.Excludes)
	for _, line := range outputs {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if includeLine(includes, line) && !excludeLine(excludes, line) {
			validOutputs = append(validOutputs, line)
		}
	}
	return validOutputs, nil
}

func Parse(ctx context.Context, cfg *Config, outputs []string) []string {
	target := cfg.Repo + "/blob/" + cfg.RepoBranch
	// replace local path to a github link
	for i, line := range outputs {
		if strings.Contains(line, cfg.RepoDir) {
			outputs[i] = strings.ReplaceAll(line, cfg.RepoDir, target)
		}
	}

	// process the example.go:7:6 -> #L7
	for i, line := range outputs {
		if strings.Contains(line, ".go:") {
			outputs[i] = strings.ReplaceAll(line, ".go:", ".go#L")
		}
	}
	return outputs
}

var divider = strings.Repeat(`=`, 100)

func PrintOutput(ctx context.Context, cfg *Config, outputs []string) {
	fmt.Printf("Run linter `%s` got %d line outputs\n", cfg.LinterCfg.LinterCommand, len(outputs))
	fmt.Println(divider)
	fmt.Printf("runner config: %+v\n", cfg)
	fmt.Println(divider)
	for _, line := range outputs {
		fmt.Println(line)
	}
	fmt.Println(divider)
	fmt.Printf("Report issue: %s/issues\n", cfg.Repo)
}

func includeLine(includes []string, line string) bool {
	if len(includes) == 0 {
		return true
	}
	for _, v := range includes {
		if strings.Contains(line, v) {
			return true
		}
	}
	return false
}

func excludeLine(excludes []string, line string) bool {
	if len(excludes) == 0 {
		return false
	}
	for _, v := range excludes {
		if strings.Contains(line, v) {
			return true
		}
	}
	return false
}

func buildIssueComment(cfg *Config, outputs []string) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Run `%s` on Repo: %s got output\n", cfg.LinterCfg.LinterCommand, cfg.Repo))
	s.WriteString("```\n")
	for _, o := range outputs {
		s.WriteString(o)
		s.WriteString("\n")
	}
	s.WriteString("```\n")
	s.WriteString(fmt.Sprintf("Report issue: %s/issues\n", cfg.Repo))
	s.WriteString(fmt.Sprintf("Github actions: %s", os.Getenv("GH_ACTION_LINK")))
	return s.String()
}

func CreateIssueComment(ctx context.Context, cfg *Config, outputs []string) error {
	body := buildIssueComment(cfg, outputs)
	cmd := exec.CommandContext(ctx, "gh", "issue", "comment",
		cfg.LinterCfg.IssueID,
		"--body", body)
	cmd.Dir = "."
	log.Printf("comment on issue #%s\n", cfg.LinterCfg.IssueID)
	return runCmd(cmd)
}
