package config

import (
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
)

func (c *Config) Validate() error {
	for name, pattern := range c.RegexPresets {
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid regex in preset %q: %w", name, err)
		}
	}

	for name, jq := range c.JQPresets {
		if _, err := gojq.Parse(jq); err != nil {
			return fmt.Errorf("invalid JQ in preset %q: %w", name, err)
		}
	}

	if c.Defaults.RegexPreset != "" {
		if _, ok := c.RegexPresets[c.Defaults.RegexPreset]; !ok {
			return fmt.Errorf("default regex_preset %q not found", c.Defaults.RegexPreset)
		}
	}
	if c.Defaults.JqPreset != "" {
		if _, ok := c.JQPresets[c.Defaults.JqPreset]; !ok {
			return fmt.Errorf("default jq_preset %q not found", c.Defaults.JqPreset)
		}
	}

	if c.Defaults.Regex != "" {
		if _, err := regexp.Compile(c.Defaults.Regex); err != nil {
			return fmt.Errorf("default regex is invalid: %w", err)
		}
	}
	if c.Defaults.JqFilter != "" {
		if _, err := gojq.Parse(c.Defaults.JqFilter); err != nil {
			return fmt.Errorf("default jq filter is invalid: %w", err)
		}
	}

	return nil
}
