package runner

import (
	"context"
	"log"
	"time"

	runner "github.com/alingse/go-linter-runner/runner/run"
)

func Run(repo string, jsonCfg string, yamlCfg string) {
	var cfg, err = runner.LoadCfg(repo, jsonCfg, yamlCfg)
	if err != nil {
		log.Fatal("load config failed: ", err)
		return
	}

	var ctx = context.Background()
	var defaultTimeout = 10 * 60 * time.Second
	ctx, cancel := context.WithTimeout(ctx, cfg.GetTimeout(defaultTimeout))
	defer cancel()

	err = runner.Prepare(ctx, cfg)
	if err != nil {
		log.Fatal("failed in prepare linter:", err)
		return
	}

	err = runner.Build(ctx, cfg)
	if err != nil {
		log.Printf("build failed and exit %+v %s", err, err.Error())
		return
	}

	outputs, err := runner.Run(ctx, cfg)
	if err != nil {
		log.Fatal("failed in run linter:", err)
		return
	}
	if len(outputs) == 0 {
		log.Println("no valid output after run")
		return
	}

	runner.Parse(ctx, cfg, outputs)
	runner.PrintOutput(ctx, cfg, outputs)
	// create comment on issue
	if cfg.LinterCfg.IssueID != "" {
		err = runner.CreateIssueComment(ctx, cfg, outputs)
		if err != nil {
			log.Fatalf("failed to CreateIssueComment err %+v \n", err)
			return
		}
	}
}
