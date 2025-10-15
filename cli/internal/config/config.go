package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiURL   string `yaml:"api_url"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func GetPath() string {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "acadule", "config.yaml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "acadule.yaml")
}

func Load() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(GetPath())
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, nil
}

func Save(cfg Config) error {
	_ = os.MkdirAll(filepath.Dir(GetPath()), 700)
	out, _ := yaml.Marshal(cfg)
	return os.WriteFile(GetPath(), out, 0600)
}
