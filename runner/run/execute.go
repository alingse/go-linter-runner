package run

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

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

func Prepare(ctx context.Context, cfg *Config) error {
	// install linter
	name, args := utils.SplitCommand(cfg.LinterCfg.InstallCommand)
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = cfg.LinterCfg.Workdir

	if err := runCmd(cmd); err != nil {
		return err
	}

	// clone repo
	cmd = exec.CommandContext(ctx, "rm", "-rf", cfg.RepoDir)
	cmd.Dir = cfg.LinterCfg.Workdir

	if err := runCmd(cmd); err != nil {
		return err
	}

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
	cfg.RepoTarget = cfg.Repo + "/blob/" + cfg.RepoBranch

	return nil
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

	for _, line := range outputs {
		// filter _test.go file
		if !cfg.LinterCfg.EnableTestfile && strings.Contains(line, testFile) {
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

//go:embed templates/issue_comment.md
var issueCommentTemplate string

type issueCommentData struct {
	GithubActionLink string
	Lines            []string
	Linter           string
	RepositoryURL    string
}

func buildIssueComment(cfg *Config, outputs []string) (string, error) {
	data := &issueCommentData{
		GithubActionLink: os.Getenv("GH_ACTION_LINK"),
		Linter:           cfg.LinterCfg.LinterCommand,
		RepositoryURL:    cfg.Repo,
	}

	for _, line := range outputs {
		text := buildIssueCommentLine(cfg, line)
		data.Lines = append(data.Lines, text)
	}

	var tpl bytes.Buffer

	tmpl, err := template.New("issue_comment").Parse(issueCommentTemplate)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func buildIssueCommentLine(cfg *Config, line string) string {
	codePath, other := buildIssueCommentLineSplit(cfg, line)
	if codePath == "" {
		return line
	}

	pathText := strings.TrimLeft(strings.ReplaceAll(codePath, cfg.RepoTarget, ""), "/:")
	codePath = cleanCodePath(codePath)
	pathText = cleanPathText(pathText)

	return fmt.Sprintf(`<a href="%s">%s</a> %s`, codePath, pathText, other)
}

func cleanCodePath(codePath string) string {
	parts := strings.Split(codePath, ":")
	if len(parts) <= 2 {
		return codePath
	}

	return strings.Join(parts[:2], ":")
}

func cleanPathText(pathText string) string {
	parts := strings.Split(pathText, ":")
	if len(parts) <= 1 {
		return pathText
	}

	return parts[0]
}

func buildIssueCommentLineSplit(cfg *Config, line string) (codePath string, other string) {
	// sytyle 1: normal linter
	if strings.Contains(line, cfg.RepoTarget) {
		return buildIssueCommentLineSplitStyle1(cfg, line)
	}
	// style 2: revive
	if strings.Contains(line, " ") {
		return buildIssueCommentLineSplitStyle2(cfg, line)
	}

	return "", line
}

func buildIssueCommentLineSplitStyle1(cfg *Config, line string) (codePath string, other string) {
	// style 1
	// /home/runner/work/go-linter-runner-example/go-linter-runner-example/rangeappendall/rangeappendslice.go:8:9: append all its data while range its
	index := strings.Index(line, cfg.RepoTarget)
	if index < 0 {
		return "", line
	}

	other = line[:index]
	tail := line[index:]
	index = strings.Index(tail, " ")

	if index < 0 {
		codePath = tail

		return strings.TrimSpace(codePath), strings.TrimSpace(other)
	}

	codePath = tail[:index]
	other += tail[index:]

	return strings.TrimSpace(codePath), strings.TrimSpace(other)
}

func buildIssueCommentLineSplitStyle2(cfg *Config, line string) (codePath string, other string) {
	// style 2: badcodes/revive/revive_modify_value.go#L17:2: suspicious assignment to a by-value method receiver (false positive?)
	parts := strings.Split(line, " ")
	others := make([]string, 0)

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.Contains(part, ".go#L") {
			if strings.HasPrefix(part, "/") {
				codePath = cfg.RepoTarget + part
			} else {
				codePath = cfg.RepoTarget + "/" + part
			}
		} else {
			others = append(others, part)
		}
	}

	if len(codePath) == 0 {
		return "", line
	}

	return codePath, strings.Join(others, " ")
}

func CreateIssueComment(ctx context.Context, cfg *Config, outputs []string) error {
	body, err := buildIssueComment(cfg, outputs)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "gh", "issue", "comment",
		cfg.LinterCfg.IssueID,
		"--body", body)
	cmd.Dir = "."

	log.Printf("comment on issue #%s\n", cfg.LinterCfg.IssueID)

	return runCmd(cmd)
}
