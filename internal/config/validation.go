package config

import (
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
)

func (c *Config) Validate() error {
	// --- Regex presets ---
	for name, pattern := range c.RegexPresets {
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid regex in preset %q: %w", name, err)
		}
	}

	// --- JQ presets ---
	for name, jq := range c.JQPresets {
		if _, err := gojq.Parse(jq); err != nil {
			return fmt.Errorf("invalid JQ in preset %q: %w", name, err)
		}
	}

	// --- Defaults ---
	d := c.Defaults

	if d.RegexPreset != "" {
		if _, ok := c.RegexPresets[d.RegexPreset]; !ok {
			return fmt.Errorf("default regex_preset %q not found", d.RegexPreset)
		}
	}

	if d.JqPreset != "" {
		if _, ok := c.JQPresets[d.JqPreset]; !ok {
			return fmt.Errorf("default jq_preset %q not found", d.JqPreset)
		}
	}

	if d.Regex != "" {
		if _, err := regexp.Compile(d.Regex); err != nil {
			return fmt.Errorf("default regex is invalid: %w", err)
		}
	}

	if d.JqFilter != "" {
		if _, err := gojq.Parse(d.JqFilter); err != nil {
			return fmt.Errorf("default jq filter is invalid: %w", err)
		}
	}

	switch d.LogFormat {
	case "", "text", "json":
	// case "", "text", "json", "csv":
		// ok
	default:
		return fmt.Errorf("log_format must be one of: text, json")
	}

	if d.Port < 1 || d.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}
