package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ActiveNotebook string `json:"active_notebook,omitempty"`
	AuthPath       string `json:"auth_path,omitempty"`
	OutputFormat   string `json:"output_format,omitempty"`
}

var (
	current *Config
	once    sync.Once
)

func Load() (*Config, error) {
	var loadErr error
	once = sync.Once{}
	once.Do(func() {
		current = &Config{}
		path := ConfigFile()
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			loadErr = err
			return
		}
		loadErr = json.Unmarshal(data, current)
	})
	return current, loadErr
}

func (c *Config) Save() error {
	if err := EnsureDir(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile(), data, 0600)
}

func GetActiveNotebook() string {
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.ActiveNotebook
}
