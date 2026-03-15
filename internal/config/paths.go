package config

import (
	"os"
	"path/filepath"
)

const appDir = ".notebooklm"

func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", appDir)
	}
	return filepath.Join(home, appDir)
}

func ConfigFile() string {
	return filepath.Join(Dir(), "config.json")
}

func AuthFile() string {
	return filepath.Join(Dir(), "storage_state.json")
}

func EnsureDir() error {
	return os.MkdirAll(Dir(), 0700)
}
