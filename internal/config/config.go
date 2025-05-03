package config

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RegexPresets map[string]string `toml:"regex_presets"`
	JQPresets    map[string]string `toml:"jq_presets"`
	Defaults     struct {
		DBFile      string `toml:"db_file"`
		Port        int    `toml:"port"`
		Launch      bool   `toml:"launch"`
		LogFormat   string `toml:"log_format"`
		RegexPreset string `toml:"regex_preset"`
		Regex       string `toml:"regex"`
		JqFilter    string `toml:"jq"`
		JqPreset    string `toml:"jq_preset"`
	} `toml:"defaults"`
}

func Load() (*Config, error) {
	path := GetConfigPath()

	return LoadFromFile(path)
}

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path := GetConfigPath()

	return SaveToFile(path, cfg)
}

func SaveToFile(path string, cfg *Config) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func GetConfigPath() string {
	path := os.Getenv("MAGIC_LOG_CONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		path = filepath.Join(home, ".magiclogrc")
	}

	return path
}
