package lib

import (
	"os"
	"path/filepath"
)

func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	config, err := filepath.Abs(filepath.Join(home, DEFAULT_CONFIG_DIRECTORY))

	return config, nil
}
