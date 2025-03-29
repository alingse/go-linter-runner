package run

import (
	"os/exec"
	"testing"
)

func TestFetchRepoInfo(t *testing.T) {
	// 检查 gh 命令是否存在
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("gh command not found, skipping test")
	}

	repoID := "go-rod/rod"
	info, err := FetchRepoInfo(repoID)
	if err != nil {
		t.Fatalf("FetchRepoInfo failed: %v", err)
	}

	// 检查返回结果不为 nil
	if info == nil {
		t.Error("RepoInfo should not be nil")
	}

	// 检查基本字段是否有值
	if info.StargazerCount <= 0 {
		t.Errorf("StargazerCount should be >= 0, got %d", info.StargazerCount)
	}
	if info.ForkCount <= 0 {
		t.Errorf("ForkCount should be >= 0, got %d", info.ForkCount)
	}
	if info.PushedAt == "" {
		t.Error("PushedAt should not be empty")
	}
	if info.UpdatedAt == "" {
		t.Error("UpdatedAt should not be empty")
	}
	t.Logf("fetch for repo %s got %#v", repoID, info)
}
