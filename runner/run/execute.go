package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/alingse/go-linter-runner/runner/utils"
)

var ErrSkipNoGoModRepo = errors.New("skip this repo for no go.mod file exists")

const (
	DiagnosticExitCode = 3
)

func runCmd(cmd *exec.Cmd) error {
	data, err := cmd.CombinedOutput()
	log.Printf("run cmd %+v got len(output)=%d and err %+v\n", cmd, len(data), err)
	fmt.Println(string(data))

	if err != nil {
		return fmt.Errorf("run %s %+v failed %w", cmd.Path, cmd.Args, err)
	}

	return nil
}

func Prepare(ctx context.Context, cfg *Config) (*RepoInfo, error) {
	var err error
	// install linter
	name, args := utils.SplitCommand(cfg.LinterCfg.InstallCommand)
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = cfg.LinterCfg.Workdir

	if err := runCmd(cmd); err != nil {
		return nil, err
	}

	// fetch repo info
	var repoInfo *RepoInfo
	if cfg.IsGithub && len(cfg.RepoID) > 0 {
		repoInfo, err = FetchRepoInfo(cfg.RepoID)
		if err != nil {
			// ignore fetch failed
			log.Printf("fetch github repo info failed %s %+v \n", cfg.RepoID, err)
		}
	}

	// skip for archived
	//if repoInfo.IsArchived {
	//	return nil
	//}
	// clone repo
	cmd = exec.CommandContext(ctx, "rm", "-rf", cfg.RepoDir)
	cmd.Dir = cfg.LinterCfg.Workdir

	if err := runCmd(cmd); err != nil {
		return nil, err
	}

	cmd = exec.CommandContext(ctx, "git", "clone", cfg.Repo)
	cmd.Dir = cfg.LinterCfg.Workdir

	if err := runCmd(cmd); err != nil {
		return nil, err
	}

	// TODO: check more deep
	// check go.mod exists
	gomodFile := path.Join(cfg.RepoDir, "go.mod")
	if !utils.IsFileExists(gomodFile) {
		return nil, ErrSkipNoGoModRepo
	}

	// run go mod download
	cmd = exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = cfg.RepoDir

	if err := runCmd(cmd); err != nil {
		return nil, err
	}

	// read default branch for repo
	cmd = exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = cfg.RepoDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git branch failed %w", err)
	}

	cfg.RepoBranch = strings.TrimSpace(string(output))
	cfg.RepoTarget = cfg.Repo + "/blob/" + cfg.RepoBranch

	return repoInfo, nil
}

func Build(ctx context.Context, cfg *Config) error {
	cmd := exec.CommandContext(ctx, "go", "build", "./...")
	cmd.Dir = cfg.RepoDir

	if err := runCmd(cmd); err != nil {
		return err
	}

	return nil
}

func Run(ctx context.Context, cfg *Config) ([]string, error) {
	name, args := utils.SplitCommand(cfg.LinterCfg.LinterCommand)
	args = append(args, "./...")
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = cfg.RepoDir

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("run cmd %+v got err %+v and exit code %+v \n", cmd, err, cmd.ProcessState.ExitCode())
	fmt.Printf("stdout:\n%s\n", stdout.String())
	fmt.Printf("stderr:\n%s\n", stderr.String())

	if err == nil {
		log.Printf("err is nil and return")

		return nil, nil
	}

	if cmd.ProcessState.ExitCode() != DiagnosticExitCode {
		log.Printf("ignore exit err %s\n", err.Error())

		return nil, nil
	}

	output := strings.TrimSpace(stderr.String())
	if len(output) == 0 {
		log.Printf("stderr output is empty, fallback to stdout")

		output = stdout.String()
		if len(output) == 0 {
			log.Printf("stdout output is still empty")

			return nil, nil
		}
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
	// replace local path to a github link
	for i, line := range outputs {
		if strings.Contains(line, cfg.RepoDir) {
			outputs[i] = strings.ReplaceAll(line, cfg.RepoDir, cfg.RepoTarget)
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

const testFile = "_test.go"

func FilterOutput(ctx context.Context, cfg *Config, outputs []string) []string {
	result := make([]string, 0, len(outputs))

	enableTestfile := utils.CastToBool(cfg.LinterCfg.EnableTestfile)
	for _, line := range outputs {
		// filter _test.go file
		if !enableTestfile && strings.Contains(line, testFile) {
			log.Println("ignore testfile output ", line)

			continue
		}

		result = append(result, line)
	}

	return result
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
