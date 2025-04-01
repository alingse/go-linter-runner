package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type RepoInfo struct {
	StargazerCount int    `json:"stargazerCount"`
	ForkCount      int    `json:"forkCount"`
	PushedAt       string `json:"pushedAt"`
	UpdatedAt      string `json:"updatedAt"`
	IsFork         bool   `json:"isFork"`
	IsEmpty        bool   `json:"isEmpty"`
	IsArchived     bool   `json:"isArchived"`
}

func FetchRepoInfo(repoID string) (*RepoInfo, error) {
	cmd := exec.Command("gh", "repo", "view", repoID, "--json", "stargazerCount,forkCount,pushedAt,updatedAt,isFork,isEmpty,isArchived")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo info: %w, stderr: %s", err, stderr.String())
	}

	var info RepoInfo
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		return nil, fmt.Errorf("failed to parse repo info: %w", err)
	}

	return &info, nil
}
