package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
	"github.com/spf13/viper"
)

func TestLoadDefaultConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".magiclogrc")

	err := os.WriteFile(path, []byte("port = 9999\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFromFile(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port != 9999 {
		t.Errorf("Expected port 9999, got %d", cfg.Port)
	}
}

func TestLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".magiclogrc")

	err := os.WriteFile(path, []byte("port = 3000\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("MAGIC_LOG_CONFIG", path)
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("viper failed to read config: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", cfg.Port)
	}
}
