package runner

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	YamlFile  string
	Repo      string
	OtherArgs RunArgs
	LinterCfg *LinterCfg
}

type LinterCfg struct {
	Name     string   `yaml:"name"`
	Linter   string   `yaml:"linter"`
	Install  string   `yaml:"install"`
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
}

type RunArgs struct {
	args map[string][]string
}

func LoadCfg() *Config {
	var cfg Config
	flag.StringVar(&cfg.YamlFile, "yaml", "", "the linter yaml config")
	flag.StringVar(&cfg.Repo, "repo", "", "the repo")
	flag.Var(&cfg.OtherArgs, "F", "the command args")
	flag.Parse()

	yamlContent, err := os.ReadFile(cfg.YamlFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		return nil
	}
	err = yaml.Unmarshal(yamlContent, &cfg.LinterCfg)
	if err != nil {
		log.Fatalf("error parsing yaml: %v", err)
		return nil
	}
	return &cfg
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
