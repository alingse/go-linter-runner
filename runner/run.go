package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var ErrSkipNoGoModRepo = errors.New("skip this repo for no go.mod file exists")

func Prepare(ctx context.Context, cfg *Config) error {
	// install linter
	args := strings.Split(cfg.LinterCfg.Install, " ")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = cfg.LinterCfg.Workdir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install linter failed %w", err)
	}
	// clone repo
	cmd = exec.CommandContext(ctx, "git", "clone", cfg.Repo)
	cmd.Dir = cfg.LinterCfg.Workdir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed %w  repo: %s", err, cfg.Repo)
	}
	// check go.mod exists
	gomodFile := path.Join(cfg.RepoDir, "go.mod")
	if !isFileExists(gomodFile) {
		return ErrSkipNoGoModRepo
	}
	// run go mod download
	cmd = exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = cfg.RepoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod download failed %w", err)
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
	cmd := exec.CommandContext(ctx, cfg.LinterCfg.Linter, "./...")
	cmd.Dir = cfg.RepoDir
	outputData, err := cmd.CombinedOutput()
	output := string(outputData)
	if err != nil && len(output) == 0 {
		return nil, err
	}
	// check includes && excludes
	outputs := strings.Split(output, "\n")
	validOutputs := make([]string, 0, len(outputs))
	for _, line := range outputs {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if includeLine(cfg, line) && !excludeLine(cfg, line) {
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
	fmt.Printf("Run %s with linter `%s` got %d line outputs\n",
		cfg.LinterCfg.Name, cfg.LinterCfg.Linter, len(outputs))
	fmt.Println(divider)
	fmt.Printf("runner config: %+v\n", cfg)
	fmt.Println(divider)
	for _, line := range outputs {
		fmt.Println(line)
	}
	fmt.Println(divider)
	fmt.Printf("Report issue: %s/issues\n", cfg.Repo)
}

func isFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func includeLine(c *Config, line string) bool {
	if len(c.LinterCfg.Includes) == 0 {
		return true
	}
	for _, v := range c.LinterCfg.Includes {
		if strings.Contains(line, v) {
			return true
		}
	}
	return false
}

func excludeLine(c *Config, line string) bool {
	if len(c.LinterCfg.Excludes) == 0 {
		return false
	}
	for _, v := range c.LinterCfg.Excludes {
		if strings.Contains(line, v) {
			return true
		}
	}
	return false
}

func CreateIssueComment(ctx context.Context, cfg *Config, outputs []string) error {
	body := fmt.Sprintf("Repo: %s\n```%s```", cfg.Repo, strings.Join(outputs, "\n"))
	cmd := exec.CommandContext(ctx, "gh", "issue", "comment",
		strconv.FormatInt(cfg.LinterCfg.Issue.ID, 10),
		"--body", body)
	cmd.Dir = "."
	data, err := cmd.CombinedOutput()
	fmt.Printf("comment on issue %d got %s and %+v \n", cfg.LinterCfg.Issue.ID, string(data), err)
	if err != nil {
		return fmt.Errorf("gh issue comment failed %w", err)
	}
	return nil
}
