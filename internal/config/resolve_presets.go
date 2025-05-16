package config

import (
	"fmt"
)

type ResolvedParsing struct {
	Regex string
	JQ    string
}

func ResolveParsingOptions(regexPreset, regexRaw, jqPreset, jqRaw string, fileCfg *Config) (ResolvedParsing, error) {
	resolvedRegex, err := ResolveRegex(regexPreset, regexRaw, fileCfg)
	if err != nil {
		return ResolvedParsing{}, err
	}

	resolvedJq, err := ResolveJqFilter(jqPreset, jqRaw, fileCfg)
	if err != nil {
		return ResolvedParsing{}, err
	}

	return ResolvedParsing{
		Regex: resolvedRegex,
		JQ:    resolvedJq,
	}, nil
}

func ResolveRegex(preset, raw string, cfg *Config) (string, error) {
	if raw != "" {
		return raw, nil
	}

	if preset != "" {
		all := GetRegexPresets(cfg)
		if regex, ok := all[preset]; ok {
			return regex, nil
		}
		return "", fmt.Errorf("unknown preset: %s", preset)
	}

	return "^(?P<message>.*)$", nil
}

func ResolveJqFilter(preset, raw string, cfg *Config) (string, error) {
	if raw != "" {
		return raw, nil
	}

	if preset != "" {
		all := GetJqPresets(cfg)
		if jqFilter, ok := all[preset]; ok {
			return jqFilter, nil
		}
		return "", fmt.Errorf("unknown preset: %s", preset)
	}

	return "", nil
}

func GetRegexPresets(cfg *Config) map[string]string {
	regex_presets := map[string]string{
		"apache": `(?P<ip>\S+) (?P<ident>\S+) (?P<user>\S+) \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<size>\d+|-)`,
	}

	// Merge in user-defined (override if names match)
	for k, v := range cfg.RegexPresets {
		regex_presets[k] = v
	}

	return regex_presets
}

func GetJqPresets(cfg *Config) map[string]string {
	jq_presets := map[string]string{}

	for k, v := range cfg.JQPresets {
		jq_presets[k] = v
	}

	return jq_presets
}
