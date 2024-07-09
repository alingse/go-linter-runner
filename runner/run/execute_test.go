package run

import (
	"fmt"
	"testing"
)

func TestBuildIssueComment(t *testing.T) {
	var outputs = []string{
		"append to slice `x` with non-zero initialized length at https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
	}
	var cfg = &Config{
		RepoTarget: "https://github.com/alingse/makezero/blob/master",
		Repo:       "https://github.com/alingse/makezero",
		LinterCfg: LinterCfg{
			LinterCommand: "makezero",
		},
	}
	body, err := buildIssueComment(cfg, outputs)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}
	t.Logf("build issue got %s \n", body)
}

func TestBuildIssueCommentLine(t *testing.T) {
	var cases = [][3]string{
		{
			"append to slice `x` with non-zero initialized length at https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
			"https://github.com/alingse/makezero/blob/master",
			"[examples/example.go#L7:6](https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6) append to slice `x` with non-zero initialized length at",
		},
		{
			"https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126:4: error-nil: use require.NoError",
			"https://github.com/eksctl-io/eksctl/blob/main",
			"[pkg/cfn/builder/managed_nodegroup_test.go#L126:4:](https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126:4:) error-nil: use require.NoError",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			line := c[0]
			cfg := &Config{RepoTarget: c[1]}
			text := buildIssueCommentLine(cfg, line)
			if text != c[2] {
				t.Errorf("expect %s but is %s", c[2], text)
			}
		})
	}
}

func TestBuildIssueCommentLineSplit(t *testing.T) {
	var cases = [][4]string{
		{
			"append to slice `x` with non-zero initialized length at https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
			"https://github.com/alingse/makezero/blob/master",
			"https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
			"append to slice `x` with non-zero initialized length at",
		},
		{
			"https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126:4: error-nil: use require.NoError",
			"https://github.com/eksctl-io/eksctl/blob/main",
			"https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126:4:",
			"error-nil: use require.NoError",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			line := c[0]
			cfg := &Config{RepoTarget: c[1]}
			code, other := buildIssueCommentLineSplit(cfg, line)
			if code != c[2] {
				t.Errorf("expect %s but is %s", c[2], code)
			}
			if other != c[3] {
				t.Errorf("expect %s but is %s", c[3], other)
			}
		})
	}
}
