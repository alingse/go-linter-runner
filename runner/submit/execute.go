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
)

func getSourceReader(source string) (io.ReadCloser, error) {
	if strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "http://") {
		return getHTTPReader(source)
	}

	if utils.IsFileExists(source) {
		return getFileReader(source)
	}

	if actionPath := os.Getenv("GITHUB_ACTION_PATH"); actionPath != "" {
		// /home/runner/work/_actions/alingse/go-linter-runner/main/source/top.txt
		filePath := path.Join(actionPath, "source", source)
		if utils.IsFileExists(filePath) {
			return getFileReader(filePath)
		}
	}

	url := "https://raw.githubusercontent.com/alingse/go-linter-runner/main/source/" + source

	return getHTTPReader(url)
}

func getFileReader(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getHTTPReader(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Get url %s failed with status %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}

func ReadSubmitRepos(source string, count int64) ([]string, error) {
	reader, err := getSourceReader(source)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	repos := make([]string, 0, int(count))
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if len(repos) >= int(count) {
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

func SumitActions(ctx context.Context, workflow string, repos []string) error {
	for i, repo := range repos {
		log.Printf("submit repo %d : %s \n", i, repo)

		err := submitRepo(ctx, workflow, repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func submitRepo(ctx context.Context, workflow string, repo string) error {
	cmd := exec.CommandContext(ctx, "gh", "workflow", "run", workflow,
		"-F", "repo_url="+repo)
	cmd.Dir = "."

	return utils.RunCmd(cmd)
}
