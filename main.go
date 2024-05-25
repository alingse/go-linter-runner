package main

import (
	"context"
	"log"
	"time"

	"github.com/alingse/go-linter-runner/runner"
)

func main() {
	var ctx = context.Background()
	var cfg = runner.LoadCfg()
	var defaultTimeout = 10 * 60 * time.Second
	var err error

	ctx, cancel := context.WithTimeout(ctx, cfg.GetTimeout(defaultTimeout))
	defer cancel()

	err = runner.Prepare(ctx, cfg)
	if err != nil {
		log.Fatal("failed in prepare linter ", err)
		return
	}
	err = runner.Run(ctx, cfg)
	if err != nil {
		log.Fatal("failed in run linter ", err)
		return
	}
}
