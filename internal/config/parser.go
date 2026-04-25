package config

import (
	"errors"
	"fmt"
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

func Parse(cfgPath string) Config {
	yml, err := os.ReadFile(cfgPath)
	if err != nil {
		fmt.Printf("Failed to read %v. Check if the file exists and has read permissions.", cfgPath)
	}

	var cfg Config

	if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
		fmt.Printf("Failed while parsing %s. Make sure that the file is properly formatted %v", cfgPath, err)
	}

	if cfg.Name == "" {
		cfg.Name = filepath.Base(filepath.Dir(cfgPath))
	}

	return cfg
}

func ResolveConfigPath(cfgPath string) (string, error) {
	if cfgPath != "" {
		exists, err := checkConfigExists(cfgPath)
		if err != nil {
			return "", fmt.Errorf("failed to check file existance")
		}
		if !exists {
			return "", fmt.Errorf("specified file does not exist %e", err)
		}

		return cfgPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory %e", err)
	}
	cwdCfgPath := filepath.Join(cwd, "/.workforest.yml")

	exists, err := checkConfigExists(cwdCfgPath)
	if err != nil {
		return "", fmt.Errorf("failed to check file existance")
	}

	if exists {
		return cwdCfgPath, nil
	}

	parentCfgPath := filepath.Join(filepath.Dir(cwd), "/.workforest.yml")

	exists, err = checkConfigExists(parentCfgPath)
	if err != nil {
		return "", fmt.Errorf("failed to check file existance")
	}

	if exists {
		return parentCfgPath, nil
	}

	return "", nil
}

func checkConfigExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}
