package config

import (
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
)

func (c *Config) Validate() []error {
	var errs []error

	// --- Regex presets ---
	for name, pattern := range c.RegexPresets {
		if _, err := regexp.Compile(pattern); err != nil {
			errs = append(errs, fmt.Errorf("regex preset %q: %v", name, err))
		}
	}

	// --- JQ presets ---
	for name, jq := range c.JQPresets {
		if _, err := gojq.Parse(jq); err != nil {
			errs = append(errs, fmt.Errorf("jq preset %q: %v", name, err))
		}
	}

	// --- Defaults ---
	d := c.Defaults

	if d.RegexPreset != "" {
		if _, ok := c.RegexPresets[d.RegexPreset]; !ok {
			errs = append(errs, fmt.Errorf("defaults.regex_preset %q not found", d.RegexPreset))
		}
	}

	if d.JqPreset != "" {
		if _, ok := c.JQPresets[d.JqPreset]; !ok {
			errs = append(errs, fmt.Errorf("defaults.jq_preset %q not found", d.JqPreset))
		}
	}

	if d.Regex != "" {
		if _, err := regexp.Compile(d.Regex); err != nil {
			errs = append(errs, fmt.Errorf("defaults.regex: %v", err))
		}
	}

	if d.JqFilter != "" {
		if _, err := gojq.Parse(d.JqFilter); err != nil {
			errs = append(errs, fmt.Errorf("defaults.jq: %v", err))
		}
	}

	switch d.LogFormat {
	case "", "text", "json":
	default:
		errs = append(errs, fmt.Errorf("defaults.log_format must be one of: text, json"))
	}

	if d.Port < 1 || d.Port > 65535 {
		errs = append(errs, fmt.Errorf("defaults.port must be between 1 and 65535"))
	}

	return errs
}
