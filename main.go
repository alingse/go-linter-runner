package main

import (
	"context"
	"log"
	"time"

	"github.com/alingse/go-linter-runner/runner"
)

func main() {
	var cfg, err = runner.LoadCfg()
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

	outputs, err := runner.Run(ctx, cfg)
	if err != nil {
		log.Fatal("failed in run linter:", err)
		return
	}
	if len(outputs) == 0 {
		log.Println("no valid output after run")
		return
	}

	err = runner.Parse(ctx, cfg, outputs)
	if err != nil {
		log.Fatal("failed parse output:", err, outputs)
		return
	}
}
