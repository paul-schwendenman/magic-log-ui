package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

type EditOptions struct {
	Editor     string
	NoBackup   bool
	NoValidate bool
}

func EditConfig(cfg EditOptions) error {
	originalPath := config.GetConfigPath()

	originalData := []byte{}
	if _, err := os.Stat(originalPath); err == nil {
		originalData, err = os.ReadFile(originalPath)
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat config: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "magiclogrc-edit-*.toml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(originalData); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	editor := cfg.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}

	cmdEdit := exec.Command(editor, tmpFile.Name())
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr
	if err := cmdEdit.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	if !cfg.NoValidate {
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
	}

	editedData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read updated config: %w", err)
	}

	writeBackup := false

	if !cfg.NoBackup {
		if _, err := os.Stat(originalPath); err == nil {
			if !bytes.Equal(originalData, editedData) {
				writeBackup = true
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat config file: %w", err)
		}
	}

	if writeBackup {
		backupPath := originalPath + ".bak"
		if err := os.WriteFile(backupPath, originalData, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Println("📦 Backup saved to:", backupPath)
	}

	if err := os.WriteFile(originalPath, editedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated config: %w", err)
	}

	fmt.Println("✅ Config updated successfully")
	return nil
}
