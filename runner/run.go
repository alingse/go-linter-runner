package runner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alingse/go-linter-runner/runner/run"
)

func Run(repo string, jsonCfg string, yamlCfg string) error {
	cfg, err := run.LoadCfg(repo, jsonCfg, yamlCfg)
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	ctx := context.Background()

	defaultTimeout := 10 * 60 * time.Second

	ctx, cancel := context.WithTimeout(ctx, cfg.GetTimeout(defaultTimeout))
	defer cancel()

	repoInfo, err := run.Prepare(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed in prepare linter: %w", err)
	}

	err = run.Build(ctx, cfg)
	if err != nil {
		log.Printf("build failed and exit %+v %s", err, err.Error())
		return nil
	}

	outputs, err := run.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed in run linter: %w", err)
	}

	if len(outputs) == 0 {
		log.Println("no valid output after run")
		return nil
	}

	// clean
	outputs = run.Parse(ctx, cfg, outputs)
	// filter
	outputs = run.FilterOutput(ctx, cfg, outputs)

	if len(outputs) == 0 {
		log.Println("no valid output after parse and filter")
		return nil
	}

	run.PrintOutput(ctx, cfg, outputs)
	// create comment on issue
	if cfg.LinterCfg.IssueID != "" {
		err = run.CreateIssueComment(ctx, cfg, repoInfo, outputs)
		if err != nil {
			return fmt.Errorf("failed to CreateIssueComment err %w", err)
		}
	}
	return nil
}
