package runner

import (
	"errors"
	"flag"
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
	OtherArgs  RunArgs
	LinterCfg  LinterCfg
	RepoURL    *url.URL
	Repo       string
	RepoDir    string
	RepoBranch string
}

type LinterCfg struct {
	Name     string   `yaml:"name"`
	Linter   string   `yaml:"linter"`
	Install  string   `yaml:"install"`
	Workdir  string   `yaml:"workdir"`
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
	Issue    IssueCfg `yaml:"issue"`
}

type IssueCfg struct {
	ID      int64 `yaml:"id"`
	Comment bool  `yaml:"comment"`
}

type RunArgs struct {
	args map[string][]string
}

func LoadCfg() (*Config, error) {
	var cfg Config
	var yamlFile string
	flag.StringVar(&yamlFile, "yaml", "", "the linter yaml config")
	flag.StringVar(&cfg.Repo, "repo", "", "the repo")
	flag.Var(&cfg.OtherArgs, "F", "the command args")
	flag.Parse()

	yamlContent, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlContent, &cfg.LinterCfg)
	if err != nil {
		return nil, err
	}
	// make workdir abs
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
	return &cfg, nil
}

func (r *RunArgs) String() string {
	return fmt.Sprintf("args %+v", r.args)
}

var ErrInvalidRunArgs = errors.New("invalid run args must be `a=b`")

func (r *RunArgs) Set(s string) error {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return ErrInvalidRunArgs
	}
	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	if parts[0] == "" || parts[1] == "" {
		return ErrInvalidRunArgs
	}
	if r.args == nil {
		r.args = make(map[string][]string)
	}
	r.args[parts[0]] = append(r.args[parts[0]], parts[1])
	return nil
}

func (r *RunArgs) Get(s string) string {
	ss, ok := r.args[s]
	if ok && len(ss) > 0 {
		return ss[0]
	}
	return ""
}

func (c *Config) GetTimeout(defaultDuration time.Duration) time.Duration {
	if v := c.OtherArgs.Get("timeout"); v != "" {
		value, _ := strconv.ParseInt(v, 10, 64)
		if value > 0 {
			return time.Duration(value) * time.Second
		}
	}
	return defaultDuration
}
