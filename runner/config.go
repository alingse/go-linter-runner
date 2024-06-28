package runner

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
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
	Includes       any    `json:"includes"`
	Excludes       any    `json:"excludes"`
	IssueID        int64  `json:"issue_id"`
	Timeout        int64  `json:"timeout"`
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
	if c.LinterCfg.Timeout > 0 {
		return time.Duration(c.LinterCfg.Timeout) * time.Second
	}
	return defaultDuration
}

func parseStringArray(s any) []string {
	var ss []string
	switch value := s.(type) {
	case string:
		if value == "" || value == "[]" {
			return ss
		}
		err := json.Unmarshal([]byte(value), &ss)
		if err != nil {
			log.Printf("parse value %s to string array failed %+v\n", value, err)
		}
		return ss
	case []any:
		if len(value) > 0 {
			data, _ := json.Marshal(ss)
			_ = json.Unmarshal(data, &ss)
			return ss
		}
	case []string:
		return value
	default:
		log.Printf("parse value %+v to string array failed\n", value)
	}
	return ss

}
