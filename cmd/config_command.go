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

	cfg, path, err := loadConfigMap()
	if err != nil {
		return err
	}

	sectionMap, ok := cfg[section].(map[string]any)
	if !ok {
		sectionMap = map[string]any{}
	}

	// 🚦 Validate and cast value for known keys
	if section == "defaults" {
		switch key {
		case "log_format":
			if value != "json" && value != "text" {
				return fmt.Errorf("log_format must be 'json' or 'text'")
			}
			sectionMap[key] = value

		case "launch":
			if value != "true" && value != "false" {
				return fmt.Errorf("launch must be 'true' or 'false'")
			}
			sectionMap[key] = (value == "true")

		case "port":
			portNum, err := strconv.Atoi(value)
			if err != nil || portNum < 1 || portNum > 65535 {
				return fmt.Errorf("port must be a number between 1 and 65535")
			}
			sectionMap[key] = portNum

		default:
			// Allow unknown default keys as strings
			sectionMap[key] = value
		}
	} else if section == "regex_presets" {
		_, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("invalid regex: %v", err)
		}
		sectionMap[key] = value
	} else {
		sectionMap[key] = value
	}

	cfg[section] = sectionMap
	return writeConfigMap(cfg, path)
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
	"parse_preset": true,
	"parse_regex":  true,
}

func normalizeKey(dotKey string) string {
	if !strings.Contains(dotKey, ".") && knownDefaults[dotKey] {
		return "defaults." + dotKey
	}
	return dotKey
}
