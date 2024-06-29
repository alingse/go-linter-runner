package run

import (
	"encoding/json"
	"errors"
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
}

type LinterCfg struct {
	Workdir        string `json:"workdir"`
	LinterCommand  string `json:"linter_command"`
	InstallCommand string `json:"install_command"`
	Includes       any    `json:"includes"`
	Excludes       any    `json:"excludes"`
	IssueID        string `json:"issue_id"`
	Timeout        string `json:"timeout"`
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
	var cfg = &Config{
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
		return nil, errors.New("install_command is empty")
	}
	if cfg.LinterCfg.LinterCommand == "" {
		return nil, errors.New("linter_command is empty")
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
