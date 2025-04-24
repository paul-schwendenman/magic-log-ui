package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Presets  map[string]string `toml:"presets"`
	Defaults struct {
		DBFile      string `toml:"db_file"`
		Port        int    `toml:"port"`
		Launch      bool   `toml:"launch"`
		LogFormat   string `toml:"log_format"`
		ParsePreset string `toml:"parse_preset"`
		ParseRegex  string `toml:"parse_regex"`
	} `toml:"defaults"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".magiclogrc")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil // no config is fine
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
