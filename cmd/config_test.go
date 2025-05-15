package cmd

import (
	"testing"
)

func TestConfigSubcommands(t *testing.T) {
	expectedSubcommands := []struct {
		use   string
		short string
	}{
		{"get <key>", "Get a config value"},
		{"set <key> <value>", "Set a config value"},
		{"unset <key>", "Unset a config value"},
		{"validate", "Validate the configuration"},
		{"edit", "Edit the raw config file in your editor"},
	}

	// Check that configCmd contains the expected subcommands.
	subcommands := configCmd.Commands()
	if len(subcommands) < len(expectedSubcommands) {
		t.Fatalf("expected at least %d subcommands, got %d", len(expectedSubcommands), len(subcommands))
	}

	for _, exp := range expectedSubcommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Use == exp.use && cmd.Short == exp.short {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("subcommand with use %q and short %q not found", exp.use, exp.short)
		}
	}
}
