package run

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LinterCfg  LinterCfg
	Repo       string
	RepoURL    *url.URL
	RepoDir    string
	RepoBranch string
	RepoTarget string
}

type LinterCfg struct {
	Workdir        string `json:"workdir"         yaml:"workdir"`
	LinterCommand  string `json:"linter_command"  yaml:"linter_command"`
	InstallCommand string `json:"install_command" yaml:"install_command"`
	Includes       any    `json:"includes"        yaml:"includes"`
	Excludes       any    `json:"excludes"        yaml:"excludes"`
	IssueID        string `json:"issue_id"        yaml:"issue_id"`
	Timeout        string `json:"timeout"         yaml:"timeout"`
	EnableTestfile any    `json:"enable_testfile" yaml:"enable_testfile"`
}

func LoadCfg(repo, jsonCfg, yamlCfg string) (*Config, error) {
	var linterCfg LinterCfg

	var err error
	if yamlCfg != "" {
		err = yaml.Unmarshal([]byte(yamlCfg), &linterCfg)
	} else {
		err = json.Unmarshal([]byte(jsonCfg), &linterCfg)
	}

	if err != nil {
		return nil, err
	}

	cfg := &Config{
		LinterCfg: linterCfg,
		Repo:      repo,
	}

	// make workdir abs
	if cfg.LinterCfg.Workdir == "" {
		cfg.LinterCfg.Workdir = "."
	}

	if !path.IsAbs(cfg.LinterCfg.Workdir) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		cfg.LinterCfg.Workdir = path.Join(cwd, cfg.LinterCfg.Workdir)
	}

	cfg.Repo = strings.TrimSuffix(cfg.Repo, "/")

	cfg.RepoURL, err = url.Parse(cfg.Repo)
	if err != nil {
		return nil, err
	}

	_, lastDir := path.Split(cfg.RepoURL.Path)
	cfg.RepoDir = path.Join(cfg.LinterCfg.Workdir, lastDir)

	// TODO: check more
	if cfg.LinterCfg.InstallCommand == "" {
		return nil, fmt.Errorf("install_command is empty %+v", cfg)
	}

	if cfg.LinterCfg.LinterCommand == "" {
		return nil, fmt.Errorf("linter_command is empty %+v", cfg)
	}

	return cfg, nil
}

func (c *Config) GetTimeout(defaultDuration time.Duration) time.Duration {
	if c.LinterCfg.Timeout != "" {
		timeout, _ := strconv.ParseInt(c.LinterCfg.Timeout, 10, 64)
		if timeout > 0 {
			return time.Duration(timeout) * time.Second
		}
	}

	return defaultDuration
}
