package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Window struct {
	Name    string
	Command string
}

type Config struct {
	Name    string
	Windows []Window
}

func Parse(cfgPath string) (Config, error) {
	yml, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
		return Config{}, err
	}

	if cfg.Name == "" {
		cfg.Name = filepath.Base(filepath.Dir(cfgPath))
	}

	return cfg, nil
}

func CheckConfigExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}
