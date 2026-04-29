package config

import "os"

func GetPath(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return cwd, nil
}

func GetDefaultConfigPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return cwd, nil
}
