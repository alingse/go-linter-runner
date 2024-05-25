package runner

import (
	"context"
	"os/exec"
)

func Prepare(ctx context.Context, cfg *Config) error {
	cmd := exec.CommandContext(ctx, cfg.LinterCfg.Install)
	return cmd.Start()
}

func Run(ctx context.Context, cfg *Config) error {
	return nil
}
