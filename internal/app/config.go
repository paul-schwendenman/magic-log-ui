package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)
import "github.com/paul-schwendenman/magic-log-ui/internal/shared"

func GetConfigValue(key string) (string, error) {
	cfg, _, err := loadConfigMap()
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
	cfg, path, err := loadConfigMap()
	if err != nil {
		return err
	}

	if !strings.Contains(dotKey, ".") {
		meta, ok := knownKeys[dotKey]
		if !ok {
			return fmt.Errorf("unsupported config key: %s", dotKey)
		}

		coerced, err := meta.Coerce(value)
		if err != nil {
			return err
		}
		cfg[dotKey] = coerced
		return writeConfigMap(cfg, path)
	}

	parts := strings.SplitN(dotKey, ".", 2)
	section, key := parts[0], parts[1]

	switch section {
	case "regex_presets":
		if _, err := shared.ValidateRegex(key)(value); err != nil {
			return err
		}
	case "jq_presets":
		if _, err := shared.ValidateJQ(key)(value); err != nil {
			return err
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
	cfg, path, err := loadConfigMap()
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

func CompleteConfigKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, _, err := loadConfigMap()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var keys []string
	for k, v := range cfg {
		switch section := v.(type) {
		case map[string]any:
			for subk := range section {
				keys = append(keys, fmt.Sprintf("%s.%s", k, subk))
			}
		default:
			keys = append(keys, k)
		}
	}

	return keys, cobra.ShellCompDirectiveNoFileComp
}

func CompleteConfigUnsetKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, _, err := loadConfigMap()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var keys []string
	for k, v := range cfg {
		if sectionMap, ok := v.(map[string]any); ok {
			for subk := range sectionMap {
				keys = append(keys, fmt.Sprintf("%s.%s", k, subk))
			}
		} else {
			keys = append(keys, k)
		}
	}

	return keys, cobra.ShellCompDirectiveNoFileComp
}

type KeyMeta struct {
	Coerce  func(string) (any, error)
	Suggest func() []string
}

var knownKeys = map[string]KeyMeta{
	"port": {
		Coerce:  shared.ParseIntInRange("port", 1, 65535),
		Suggest: nil,
	},
	"launch": {
		Coerce:  shared.ParseBool("launch"),
		Suggest: shared.SuggestBool,
	},
	"has_csv_header": {
		Coerce:  shared.ParseBool("has_csv_header"),
		Suggest: shared.SuggestBool,
	},
	"log_format": {
		Coerce:  shared.ParseEnum("log_format", "json", "text"),
		Suggest: func() []string { return []string{"json", "text", "csv"} },
	},
	"regex": {
		Coerce:  shared.ValidateRegex("regex"),
		Suggest: nil,
	},
	"regex_preset": {
		Coerce:  shared.StringPassThrough("regex_preset"),
		Suggest: func() []string { return getKeysFromSection("regex_presets") },
	},
	"jq": {
		Coerce:  shared.ValidateJQ("jq"),
		Suggest: nil,
	},
	"jq_preset": {
		Coerce:  shared.StringPassThrough("jq_preset"),
		Suggest: func() []string { return getKeysFromSection("jq_presets") },
	},
	"csv_fields": {
		Coerce:  shared.StringPassThrough("csv_fields"),
		Suggest: nil,
	},
}

var knownSections = []string{
	"regex_presets",
	"jq_presets",
}

func CompleteKnownConfigKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		var keys []string

		// Top-level keys
		for k := range knownKeys {
			keys = append(keys, k)
		}

		// Section stubs
		for _, s := range knownSections {
			keys = append(keys, s+".")
		}

		return keys, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 1 {
		return suggestValuesForKey(args[0]), cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func completeConfigKeyValues(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// completing first arg: the key
		return getAllConfigKeys(), cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 1 {
		key := args[0]
		return suggestValuesForKey(key), cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func getAllConfigKeys() []string {
	cfg, _, err := loadConfigMap()
	if err != nil {
		return nil
	}

	var keys []string
	for k, v := range cfg {
		if section, ok := v.(map[string]any); ok {
			for subk := range section {
				keys = append(keys, fmt.Sprintf("%s.%s", k, subk))
			}
		} else {
			keys = append(keys, k)
		}
	}
	return keys
}

func suggestValuesForKey(key string) []string {
	if meta, ok := knownKeys[key]; ok && meta.Suggest != nil {
		return meta.Suggest()
	}
	return nil
}

func getKeysFromSection(section string) []string {
	cfg, _, err := loadConfigMap()
	if err != nil {
		return nil
	}
	sec, ok := cfg[section].(map[string]any)
	if !ok {
		return nil
	}
	var keys []string
	for k := range sec {
		keys = append(keys, k)
	}
	return keys
}
