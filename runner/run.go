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
	if err := cmd.Start(); err != nil {
		return err
	}
	// clone repo
	cmd = exec.CommandContext(ctx, "git", "clone", cfg.Repo)
	cmd.Dir = cfg.LinterCfg.Workdir
	if err := cmd.Start(); err != nil {
		return err
	}
	// check go.mod exists
	if !isFileExists(path.Join(cfg.RepoDir, "go.mod")) {
		return ErrSkipNoGoModRepo
	}
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

func Parse(ctx context.Context, cfg *Config, outputs []string) error {
	for i, line := range outputs {
		if strings.Contains(line, cfg.RepoDir) {
			outputs[i] = strings.ReplaceAll(line, cfg.RepoDir, cfg.Repo)
		}
	}
	fmt.Println(outputs)
	return nil
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
