package config

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DBFile       string            `toml:"db_file" json:"db_file,omitempty"`
	Port         int               `toml:"port" json:"port,omitempty"`
	Launch       bool              `toml:"launch" json:"launch,omitempty"`
	LogFormat    string            `toml:"log_format" json:"log_format,omitempty"`
	RegexPreset  string            `toml:"regex_preset" json:"regex_preset,omitempty"`
	Regex        string            `toml:"regex" json:"regex,omitempty"`
	JqFilter     string            `toml:"jq" json:"jq,omitempty"`
	JqPreset     string            `toml:"jq_preset" json:"jq_preset,omitempty"`
	CSVFields    string            `toml:"csv_fields" json:"csv_fields,omitempty"`
	HasCSVHeader bool              `toml:"has_csv_header" json:"has_csv_header,omitempty"`

	RegexPresets map[string]string `toml:"regex_presets" json:"regex_presets,omitempty"`
	JQPresets    map[string]string `toml:"jq_presets" json:"jq_presets,omitempty"`
}

func Load() (*Config, error) {
	path := GetConfigPath()

	return LoadFromFile(path)
}

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				RegexPresets: map[string]string{},
				JQPresets:    map[string]string{},
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.RegexPresets == nil {
		cfg.RegexPresets = make(map[string]string)
	}
	if cfg.JQPresets == nil {
		cfg.JQPresets = make(map[string]string)
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
