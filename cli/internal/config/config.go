package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiURL   string `yaml:"api_url"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

var (
	configPath *string = nil
)

func init() {
	resolvedConfigPath, err := resolveConfigPath()
	if err != nil {
		slog.Error("Failed to get config file path", slog.Any("error", err))
		panic("Failed to resolve config path")
	}
	configPath = resolvedConfigPath
}

// getConfigFolder returns config folder path.
func getConfigFolder() (*string, error) {
	// get user config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// returns config dir for acadule
	configDir := filepath.Join(userConfigDir, "acadule")
	return &configDir, nil
}

// resolveConfigPath returns resolved config file path
func resolveConfigPath() (*string, error) {
	// get config dir
	configDir, err := getConfigFolder()
	if err != nil {
		return nil, err
	}

	// returns config file
	configFile := filepath.Join(*configDir, "config.yaml")
	return &configFile, nil
}

// GetConfigPath returns config file path
func GetConfigPath() string {
	if configPath == nil {
		panic("Config path not found.")
	}
	return *configPath
}

func Load() (Config, error) {
	var cfg Config = Config{}
	data, err := os.ReadFile(GetConfigPath())
	if errors.Is(err, os.ErrNotExist) {
		// return default
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, nil
}

func Save(cfg Config) error {
	_ = os.MkdirAll(filepath.Dir(GetConfigPath()), 0700)
	out, _ := yaml.Marshal(cfg)
	return os.WriteFile(GetConfigPath(), out, 0600)
}
