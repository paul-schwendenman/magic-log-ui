package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

func TestLoadDefaultConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".magiclogrc")

	err := os.WriteFile(path, []byte("[defaults]\nport = 3000\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFromFile(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Defaults.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", cfg.Defaults.Port)
	}
}
