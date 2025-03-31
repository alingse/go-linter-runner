package run

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/issue_comment.md
var issueCommentTemplate string

type issueCommentData struct {
	GithubActionLink string
	Lines            []string
	Linter           string
	RepositoryURL    string
	RepoInfo         *RepoInfo
}

var templateFuncs = template.FuncMap{
	"formatCount": formatCount,
	"isOldDate":   isOldDate,
	"formatDate":  formatDate,
	"repoName":    repoName,
	"yearsSince":  yearsSince,
}

func buildIssueComment(cfg *Config, repoInfo *RepoInfo, outputs []string) (string, error) {
	data := &issueCommentData{
		GithubActionLink: os.Getenv("GH_ACTION_LINK"),
		Linter:           cfg.LinterCfg.LinterCommand,
		RepositoryURL:    cfg.Repo,
		RepoInfo:         repoInfo,
	}

	for _, line := range outputs {
		text := buildIssueCommentLine(cfg, line)
		data.Lines = append(data.Lines, text)
	}

	var tpl bytes.Buffer

	tmpl, err := template.New("issue_comment").Funcs(templateFuncs).Parse(issueCommentTemplate)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}

	comment := tpl.String()
	comment = strings.TrimSpace(comment)
	return comment, nil
}

func formatCount(count int) string {
	if count >= 1000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func isOldDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return false
	}
	return time.Since(t) > 365*24*time.Hour
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format(time.DateTime)
}

func yearsSince(dateStr string) int {
	if dateStr == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return 0
	}
	return int(time.Since(t).Hours() / 24 / 365)
}

func repoName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return url
	}
	return parts[len(parts)-1]
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

func CreateIssueComment(ctx context.Context, cfg *Config, repoInfo *RepoInfo, outputs []string) error {
	body, err := buildIssueComment(cfg, repoInfo, outputs)
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
