package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/itchyny/gojq"
)

func GetConfigValue(key string) (string, error) {
	cfg, _, err := LoadConfigMap()
	if err != nil {
		return "", err
	}

	if val, ok := cfg[key]; ok {
		return fmt.Sprintf("%v", val), nil
	}

	if strings.Contains(key, ".") {
		parts := strings.SplitN(key, ".", 2)
		section, subkey := parts[0], parts[1]

		sectionMap, ok := cfg[section].(map[string]any)
		if !ok {
			return "", fmt.Errorf("no such section: %s", section)
		}
		val, ok := sectionMap[subkey]
		if !ok {
			return "", fmt.Errorf("no such key: %s in section %s", subkey, section)
		}
		return fmt.Sprintf("%v", val), nil
	}

	return "", fmt.Errorf("key not found: %s", key)
}

func SetConfigValue(dotKey, value string) error {
	cfg, path, err := LoadConfigMap()
	if err != nil {
		return err
	}

	if !strings.Contains(dotKey, ".") {
		switch dotKey {
		case "port":
			p, err := strconv.Atoi(value)
			if err != nil || p < 1 || p > 65535 {
				return fmt.Errorf("invalid port: %s", value)
			}
			cfg[dotKey] = p

		case "launch":
			cfg[dotKey] = (value == "true")

		case "log_format":
			if value != "text" && value != "json" {
				return fmt.Errorf("log_format must be 'json' or 'text'")
			}
			cfg[dotKey] = value

		case "regex", "jq", "regex_preset", "jq_preset", "csv_fields":
			cfg[dotKey] = value

		case "has_csv_header":
			cfg[dotKey] = (value == "true")

		default:
			return fmt.Errorf("unsupported config key: %s", dotKey)
		}

		return writeConfigMap(cfg, path)
	}

	parts := strings.SplitN(dotKey, ".", 2)
	section, key := parts[0], parts[1]

	switch section {
	case "regex_presets":
		if _, err := regexp.Compile(value); err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
	case "jq_presets":
		if _, err := gojq.Parse(value); err != nil {
			return fmt.Errorf("invalid jq expression: %w", err)
		}
	default:
		return fmt.Errorf("unknown section: %s", section)
	}

	sectionMap, ok := cfg[section].(map[string]any)
	if !ok {
		sectionMap = map[string]any{}
	}
	sectionMap[key] = value
	cfg[section] = sectionMap

	return writeConfigMap(cfg, path)
}


func UnsetConfigValue(dotKey string) error {
	cfg, path, err := LoadConfigMap()
	if err != nil {
		return err
	}

	if !strings.Contains(dotKey, ".") {
		if _, exists := cfg[dotKey]; !exists {
			return fmt.Errorf("no such key: %s", dotKey)
		}
		delete(cfg, dotKey)
		return writeConfigMap(cfg, path)
	}

	parts := strings.SplitN(dotKey, ".", 2)
	section, key := parts[0], parts[1]

	sectionMap, ok := cfg[section].(map[string]any)
	if !ok {
		return fmt.Errorf("no such section: %s", section)
	}

	if _, exists := sectionMap[key]; !exists {
		return fmt.Errorf("no such key: %s in section %s", key, section)
	}

	delete(sectionMap, key)
	if len(sectionMap) == 0 {
		delete(cfg, section)
	} else {
		cfg[section] = sectionMap
	}

	return writeConfigMap(cfg, path)
}


func LoadConfigMap() (map[string]any, string, error) {
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
