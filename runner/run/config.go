package run

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	LinterCfg  LinterCfg
	RepoURL    *url.URL
	RepoDir    string
	RepoBranch string
}

type LinterCfg struct {
	Workdir        string `json:"workdir"`
	LinterCommand  string `json:"linter_command"`
	InstallCommand string `json:"install_command"`
	RepoURL        string `json:"repo_url"`
	Includes       string `json:"includes"`
	Excludes       string `json:"excludes"`
	IssueID        string `json:"issue_id"`
	Timeout        string `json:"timeout"`
}

func LoadCfg(arg string) (*Config, error) {
	var linterCfg LinterCfg
	err := json.Unmarshal([]byte(arg), &linterCfg)
	if err != nil {
		return nil, err
	}
	var cfg = &Config{
		LinterCfg: linterCfg,
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

	cfg.LinterCfg.RepoURL = strings.TrimSuffix(cfg.LinterCfg.RepoURL, "/")
	cfg.RepoURL, err = url.Parse(cfg.LinterCfg.RepoURL)
	if err != nil {
		return nil, err
	}

	_, lastDir := path.Split(cfg.RepoURL.Path)
	cfg.RepoDir = path.Join(cfg.LinterCfg.Workdir, lastDir)
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

func parseStringArray(s string) []string {
	var ss []string
	if s == "" || s == "[]" {
		return ss
	}
	err := json.Unmarshal([]byte(s), &ss)
	if err != nil {
		log.Printf("parse value %s to string array failed %+v\n", s, err)
	}
	return ss
}
