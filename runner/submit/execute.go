package submit

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/alingse/go-linter-runner/runner/utils"
	"golang.org/x/time/rate"
)

func getSourceReader(ctx context.Context, cfg *SubmitConfig) (io.ReadCloser, error) {
	if strings.HasPrefix(cfg.Source, "https://") || strings.HasPrefix(cfg.Source, "http://") {
		return getHTTPReader(ctx, cfg.Source)
	}

	if utils.IsFileExists(cfg.Source) {
		return getFileReader(cfg.Source)
	}

	if actionPath := os.Getenv("GITHUB_ACTION_PATH"); actionPath != "" {
		// /home/runner/work/_actions/alingse/go-linter-runner/main/source/top.txt
		filePath := path.Join(actionPath, "source", cfg.Source)
		if utils.IsFileExists(filePath) {
			return getFileReader(filePath)
		}
	}

	url := "https://raw.githubusercontent.com/alingse/go-linter-runner/main/source/" + cfg.Source

	return getHTTPReader(ctx, url)
}

func getFileReader(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getHTTPReader(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Get url %s failed with status %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}

func ReadSubmitRepos(ctx context.Context, cfg *SubmitConfig) ([]string, error) {
	reader, err := getSourceReader(ctx, cfg)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	repos := make([]string, 0, cfg.RepoCount)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if len(repos) >= cfg.RepoCount {
			return repos, nil
		}

		line := scanner.Text()

		repo, err := url.Parse(line)
		if err != nil {
			return nil, err
		}

		repos = append(repos, repo.String())
	}

	return repos, nil
}

func SumitActions(ctx context.Context, cfg *SubmitConfig, repos []string) error {
	limiter := rate.NewLimiter(rate.Limit(cfg.RateLimit), 1)

	for i, repo := range repos {
		log.Printf("submit repo %d : %s \n", i, repo)

		if err := limiter.Wait(ctx); err != nil {
			return err
		}

		err := submitRepo(ctx, cfg, repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func submitRepo(ctx context.Context, cfg *SubmitConfig, repo string) error {
	args := []string{
		"workflow",
		"run",
		cfg.Workflow,
	}
	if cfg.WorkflowRef != "" {
		args = append(args, "-r", cfg.WorkflowRef)
	}
	args = append(args, "-F", "repo_url="+repo)

	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Dir = "."
	return utils.RunCmd(cmd)
}
