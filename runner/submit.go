package runner

import (
	"context"

	"github.com/alingse/go-linter-runner/runner/submit"
)

const (
	DefaultSource string = `top.txt`
	DefaultCount  int64  = 2000
)

func Submit(ctx context.Context, cfg *submit.SubmitConfig) error {
	repos, err := submit.ReadSubmitRepos(ctx, cfg)
	if err != nil {
		return err
	}
	return submit.SumitActions(ctx, cfg, repos)
}
