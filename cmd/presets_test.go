package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestPresetsCmd(t *testing.T) {
	// Create a temporary config file with minimal config content.
	tmpFile, err := os.CreateTemp("", "testconfig-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
db_file = "dummy.db"
port = 8080
launch = false
log_format = "plain"

[regex_presets]
apache = "(?P<ip>\\d+\\.\\d+\\.\\d+\\.\\d+)"

[jq_presets]
filter1 = ".data"
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Set viper to use the temporary config file.
	viper.SetConfigFile(tmpFile.Name())
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	// Capture stdout.
	oldOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Execute the presets command.
	presetsCmd.Run(nil, []string{})

	// Restore stdout and capture the output.
	w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	os.Stdout = oldOut

	output := buf.String()
	if !strings.Contains(output, "apache") {
		t.Errorf("expected output to contain 'apache' preset, got %s", output)
	}
	if !strings.Contains(output, "filter1") {
		t.Errorf("expected output to contain 'filter1' preset, got %s", output)
	}
}
