package run

import (
	"strings"
	"testing"
	"time"
)

func TestBuildIssueCommentWithRepoInfo(t *testing.T) {
	outputs := []string{
		"example output line 1",
		"example output line 2",
	}

	baseCfg := &Config{
		RepoTarget: "https://github.com/example/repo/blob/main",
		Repo:       "https://github.com/example/repo",
		LinterCfg: LinterCfg{
			LinterCommand: "example-linter",
		},
	}

	tests := []struct {
		name     string
		cfg      *Config
		repoInfo *RepoInfo
		wantErr  bool
	}{
		{
			name:     "nil repoInfo",
			cfg:      baseCfg,
			repoInfo: nil,
			wantErr:  false,
		},
		{
			name: "normal repo",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 5792,
				ForkCount:      369,
				PushedAt:       "2024-12-07T14:51:24Z",
				UpdatedAt:      "2025-03-30T01:26:21Z",
				IsFork:         false,
				IsEmpty:        false,
				IsArchived:     false,
			},
			wantErr: false,
		},
		{
			name: "archived repo",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 1200,
				ForkCount:      50,
				PushedAt:       "2024-01-01T00:00:00Z",
				UpdatedAt:      "2024-01-01T00:00:00Z",
				IsFork:         false,
				IsEmpty:        false,
				IsArchived:     true,
			},
			wantErr: false,
		},
		{
			name: "old repo",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 500,
				ForkCount:      20,
				PushedAt:       "2020-01-01T00:00:00Z",
				UpdatedAt:      "2020-01-01T00:00:00Z",
				IsFork:         false,
				IsEmpty:        false,
				IsArchived:     false,
			},
			wantErr: false,
		},
		{
			name: "empty repo",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 0,
				ForkCount:      0,
				PushedAt:       time.Now().Format(time.RFC3339),
				UpdatedAt:      time.Now().Format(time.RFC3339),
				IsFork:         false,
				IsEmpty:        true,
				IsArchived:     false,
			},
			wantErr: false,
		},
		{
			name: "fork repo",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 100,
				ForkCount:      5,
				PushedAt:       time.Now().Format(time.RFC3339),
				UpdatedAt:      time.Now().Format(time.RFC3339),
				IsFork:         true,
				IsEmpty:        false,
				IsArchived:     false,
			},
			wantErr: false,
		},
		{
			name: "large numbers",
			cfg:  baseCfg,
			repoInfo: &RepoInfo{
				StargazerCount: 1500,
				ForkCount:      1200,
				PushedAt:       time.Now().Format(time.RFC3339),
				UpdatedAt:      time.Now().Format(time.RFC3339),
				IsFork:         false,
				IsEmpty:        false,
				IsArchived:     false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GH_ACTION_LINK", "https://github.com/example/action")

			got, err := buildIssueComment(tt.cfg, tt.repoInfo, outputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildIssueComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Basic content checks
			if !strings.Contains(got, tt.cfg.LinterCfg.LinterCommand) {
				t.Error("output should contain linter command")
			}
			if !strings.Contains(got, tt.cfg.Repo) {
				t.Error("output should contain repo URL")
			}

			// Check repo info display
			if tt.repoInfo != nil {
				if tt.repoInfo.IsArchived && !strings.Contains(got, "This repository is archived") {
					t.Error("archived repo should show warning")
				}
				if tt.repoInfo.IsFork && !strings.Contains(got, "This is a fork repository") {
					t.Error("fork repo should show indicator")
				}
				if tt.repoInfo.IsEmpty && !strings.Contains(got, "This is an empty repository") {
					t.Error("empty repo should show indicator")
				}
				if tt.repoInfo.StargazerCount >= 1000 && !strings.Contains(got, "k") {
					t.Error("large star count should be formatted")
				}
				if tt.repoInfo.ForkCount >= 1000 && !strings.Contains(got, "k") {
					t.Error("large fork count should be formatted")
				}
			} else {
				if !strings.Contains(got, "Failed to get repository details") {
					t.Error("nil repoInfo should show error message")
				}
			}

			// Check output lines
			for _, line := range outputs {
				if !strings.Contains(got, line) {
					t.Errorf("output should contain line: %s", line)
				}
			}
		})
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want string
	}{
		{
			name: "less than 1000",
			arg:  999,
			want: "999",
		},
		{
			name: "exactly 1000",
			arg:  1000,
			want: "1.0k",
		},
		{
			name: "greater than 1000",
			arg:  1500,
			want: "1.5k",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatCount(tt.arg); got != tt.want {
				t.Errorf("formatCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOldDate(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "empty date",
			arg:  "",
			want: false,
		},
		{
			name: "recent date",
			arg:  time.Now().Format(time.RFC3339),
			want: false,
		},
		{
			name: "old date",
			arg:  time.Now().AddDate(-2, 0, 0).Format(time.RFC3339),
			want: true,
		},
		{
			name: "invalid date",
			arg:  "invalid-date",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOldDate(tt.arg); got != tt.want {
				t.Errorf("isOldDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "valid date",
			arg:  now.Format(time.RFC3339),
			want: now.Format(time.DateTime),
		},
		{
			name: "invalid date",
			arg:  "invalid-date",
			want: "invalid-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDate(tt.arg); got != tt.want {
				t.Errorf("formatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoName(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "normal repo URL",
			arg:  "https://github.com/example/repo",
			want: "repo",
		},
		{
			name: "short URL",
			arg:  "https://github.com/repo",
			want: "repo",
		},
		{
			name: "invalid URL",
			arg:  "not-a-url",
			want: "not-a-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repoName(tt.arg); got != tt.want {
				t.Errorf("repoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
