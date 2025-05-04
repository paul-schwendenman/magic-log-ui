package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func handleConfigCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: magic-log config [get|set|unset] <key> [value]")
		os.Exit(1)
	}

	cmd := args[0]
	key := args[1]

	switch cmd {
	case "get":
		val, err := GetConfigValue(key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
		if val == "" {
			os.Exit(1)
		}
		fmt.Println(val)

	case "set":
		if len(args) < 3 {
			fmt.Println("Usage: magic-log config set <key> <value>")
			os.Exit(1)
		}
		val := args[2]
		if err := SetConfigValue(key, val); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}

	case "unset":
		if err := UnsetConfigValue(key); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Usage: magic-log config [get|set|unset] <key> [value]")
		os.Exit(1)
	}
}

func loadConfigMap() (map[string]any, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(home, ".magiclogrc")

	config := map[string]any{}
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, "", err
		}
		if err := toml.Unmarshal(data, &config); err != nil {
			return nil, "", err
		}
	}
	return config, path, nil
}

func writeConfigMap(cfg map[string]any, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(cfg)
}

func GetConfigValue(dotKey string) (string, error) {
	dotKey = normalizeKey(dotKey)

	cfg, _, err := loadConfigMap()
	if err != nil {
		return "", err
	}

	parts := strings.Split(dotKey, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid key format: use section.key")
	}
	section, key := parts[0], parts[1]

	sectionMap, ok := cfg[section].(map[string]any)
	if !ok {
		return "", nil
	}
	val, ok := sectionMap[key]
	if !ok {
		return "", nil
	}

	return fmt.Sprintf("%v", val), nil
}

func SetConfigValue(dotKey, value string) error {
	dotKey = normalizeKey(dotKey)

	parts := strings.Split(dotKey, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format: use section.key")
	}
	section, key := parts[0], parts[1]

	// 🔄 Load current typed config
	typedCfg, err := config.Load()
	if err != nil {
		return err
	}

	// 🧪 Apply and validate in memory
	switch section {
	case "defaults":
		switch key {
		case "log_format":
			if value != "text" && value != "json" {
				return fmt.Errorf("log_format must be 'json' or 'text'")
			}
			typedCfg.Defaults.LogFormat = value
		case "port":
			p, err := strconv.Atoi(value)
			if err != nil || p < 1 || p > 65535 {
				return fmt.Errorf("invalid port: %s", value)
			}
			typedCfg.Defaults.Port = p
		case "launch":
			typedCfg.Defaults.Launch = (value == "true")
		case "regex_preset":
			typedCfg.Defaults.RegexPreset = value
		case "jq_preset":
			typedCfg.Defaults.JqPreset = value
		default:
			return fmt.Errorf("unsupported default key: %s", key)
		}
	case "regex_presets":
		if _, err := regexp.Compile(value); err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
		typedCfg.RegexPresets[key] = value

	case "jq_presets":
		if _, err := gojq.Parse(value); err != nil {
			return fmt.Errorf("invalid jq expression: %w", err)
		}
		typedCfg.JQPresets[key] = value

	default:
		return fmt.Errorf("unknown section: %s", section)
	}

	// ✅ Validate full config
	if err := typedCfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 💾 Save back to file
	path := configPath()
	return config.SaveToFile(path, typedCfg)
}

func UnsetConfigValue(dotKey string) error {
	cfg, path, err := loadConfigMap()
	if err != nil {
		return err
	}

	dotKey = normalizeKey(dotKey)
	parts := strings.Split(dotKey, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format: use section.key")
	}
	section, key := parts[0], parts[1]

	sectionMap, ok := cfg[section].(map[string]any)
	if !ok {
		return fmt.Errorf("no such section: %s", section)
	}
	delete(sectionMap, key)
	if len(sectionMap) == 0 {
		delete(cfg, section)
	} else {
		cfg[section] = sectionMap
	}

	return writeConfigMap(cfg, path)
}

var knownDefaults = map[string]bool{
	"log_format":   true,
	"port":         true,
	"launch":       true,
	"db_file":      true,
	"regex_preset": true,
	"regex":        true,
	"jq":           true,
	"jq_preset":    true,
}

func normalizeKey(dotKey string) string {
	if !strings.Contains(dotKey, ".") && knownDefaults[dotKey] {
		return "defaults." + dotKey
	}
	return dotKey
}
