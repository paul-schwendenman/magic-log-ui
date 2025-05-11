package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

func EditConfig() error {
	originalPath := config.GetConfigPath()

	// Load original contents
	originalData, err := os.ReadFile(originalPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "magiclogrc-edit-*.toml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // cleanup temp file

	if _, err := tmpFile.Write(originalData); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	// Open with editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // fallback
	}
	cmdEdit := exec.Command(editor, tmpFile.Name())
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr
	if err := cmdEdit.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	// Re-validate edited config
	editedCfg, err := config.LoadFromFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("❌ Failed to parse edited config: %w", err)
	}
	if errs := editedCfg.Validate(); len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "❌ Edited config is invalid:")
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "   -", e)
		}
		return fmt.Errorf("aborting due to invalid config")
	}

	// Move it back to original path
	if err := os.WriteFile(originalPath, originalData, 0644); err != nil {
		return fmt.Errorf("failed to write updated config: %w", err)
	}

	fmt.Println("✅ Config updated successfully")
	return nil
}
