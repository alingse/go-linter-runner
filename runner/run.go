package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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
		return fmt.Errorf("git clon failed %w", err)
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
	if err != nil {
		return nil, err
	}
	// check includes && excludes
	output := string(outputData)
	outputs := strings.Split(output, "\n")
	validOutputs := make([]string, 0, len(outputs))
	for _, line := range outputs {
		if includeLine(cfg, line) && !excludeLine(cfg, line) {
			validOutputs = append(validOutputs, output)
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

func Print(ctx context.Context, cfg *Config, outputs []string) {
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
