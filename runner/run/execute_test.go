package run

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestBuildIssueComment(t *testing.T) {
	outputs := []string{
		"append to slice `x` with non-zero initialized length at https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
	}

	cfg := &Config{
		RepoTarget: "https://github.com/alingse/makezero/blob/master",
		Repo:       "https://github.com/alingse/makezero",
		LinterCfg: LinterCfg{
			LinterCommand: "makezero",
		},
	}

	body, err := buildIssueComment(cfg, nil, outputs)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}

	t.Logf("build issue got %s \n", body)
}

func TestBuildIssueCommentLine(t *testing.T) {
	cases := [][3]string{
		{
			"append to slice `x` with non-zero initialized length at https://github.com/alingse/makezero/blob/master/examples/example.go#L7:6",
			"https://github.com/alingse/makezero/blob/master",
			"<a href=\"https://github.com/alingse/makezero/blob/master/examples/example.go#L7\">examples/example.go#L7</a> append to slice `x` with non-zero initialized length at",
		},
		{
			"https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126:4: error-nil: use require.NoError",
			"https://github.com/eksctl-io/eksctl/blob/main",
			"<a href=\"https://github.com/eksctl-io/eksctl/blob/main/pkg/cfn/builder/managed_nodegroup_test.go#L126\">pkg/cfn/builder/managed_nodegroup_test.go#L126</a> error-nil: use require.NoError",
		},
		{
			"badcodes/revive/revive_modify_value.go#L17:2: suspicious assignment to a by-value method receiver (false positive?)",
			"https://github.com/alingse/go-linter-runner-example/blob/main",
			"<a href=\"https://github.com/alingse/go-linter-runner-example/blob/main/badcodes/revive/revive_modify_value.go#L17\">badcodes/revive/revive_modify_value.go#L17</a> suspicious assignment to a by-value method receiver (false positive?)",
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
	cases := [][4]string{
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

func TestOutput(t *testing.T) {
	ctx := context.Background()

	output := `badcodes/revive/revive_modify_value.go:17:2: suspicious assignment to a by-value method receiver (false positive?)
badcodes/revive/revive_modify_value.go:22:2: suspicious assignment to a by-value method receiver (false positive?)`

	outputs := strings.Split(output, "\n")

	repo := "https://github.com/alingse/go-linter-runner-example"
	repoURL, _ := url.Parse(repo)
	cfg := &Config{
		RepoTarget: "https://github.com/alingse/go-linter-runner-example/blob/main",
		RepoDir:    "/home/",
		RepoURL:    repoURL,
		Repo:       repo,
	}

	t.Setenv("GH_ACTION_LINK", "https://github.com/xxx")

	outputs = Parse(ctx, cfg, outputs)

	body, err := buildIssueComment(cfg, nil, outputs)
	if err != nil {
		t.Errorf("err should be nil but got %+v", err)
	}

	var expected = `Run ` + "``" + ` on Repo: https://github.com/alingse/go-linter-runner-example


<div style="border: 1px solid #e1e4e8; border-radius: 6px; padding: 16px; margin: 16px 0;">
  <h3 style="margin-top: 0;">go-linter-runner-example</h3>
  <div>⚠️ Failed to get repository details</div>
</div>


Got total 2 lines output in action: https://github.com/xxx

<details open>
<summary>Click to expand details</summary>
<ol>
<li><a href="https://github.com/alingse/go-linter-runner-example/blob/main/badcodes/revive/revive_modify_value.go#L17">badcodes/revive/revive_modify_value.go#L17</a> suspicious assignment to a by-value method receiver (false positive?)</li>
<li><a href="https://github.com/alingse/go-linter-runner-example/blob/main/badcodes/revive/revive_modify_value.go#L22">badcodes/revive/revive_modify_value.go#L22</a> suspicious assignment to a by-value method receiver (false positive?)</li></ol>
</details>

Report issue: https://github.com/alingse/go-linter-runner-example/issues`

	if body != expected {
		t.Errorf("body %s is not expect %s", body, expected)
	}
}
