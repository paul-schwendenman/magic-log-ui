package config

import (
	"flag"
	"fmt"
	"os"
)

type FinalConfig struct {
	DBFile       string
	Port         int
	Launch       bool
	Echo         bool
	LogFormat    string
	RegexPreset  string
	Regex        string
	JqFilter     string
	JqPreset     string
	CSVFieldsStr string
	HasCSVHeader bool
	AutoAnalyze  bool
	ShowVersion  bool
	ListPresets  bool
}

func ParseArgsAndConfig() (*FinalConfig, *Config, error) {
	var (
		dbFile        = flag.String("db-file", "", "Path to a DuckDB database file.")
		port          = flag.Int("port", 3000, "Port to serve the web UI on.")
		launch        = flag.Bool("launch", false, "Open the UI in a browser.")
		echo          = flag.Bool("echo", false, "Echo parsed stdin input to stdout")
		noAutoAnalyze = flag.Bool("no-auto-analyze", false, "Disable automatic ANALYZE of logs table")
		logFormat     = flag.String("log-format", "json", "Log format: json, csv or plain text.")
		parseRegex    = flag.String("regex", "", "Custom regex to parse logs. Use with \"text\" format")
		regexPreset   = flag.String("regex-preset", "", "Regex preset to use.")
		jqFilter      = flag.String("jq", "", "A jq expression to apply to parsed logs")
		jqPreset      = flag.String("jq-preset", "", "Regex preset to use.")
		csvFieldsStr  = flag.String("csv-fields", "", "Comma-separated field names for CSV logs (used with --log-format=csv)")
		hasCSVHeader  = flag.Bool("has-csv-header", true, "Indicates if CSV logs include a header row")
		listPresets   = flag.Bool("list-presets", false, "List available regex and jq presets and exit.")
		showVersion   = flag.Bool("version", false, "Print version and exit.")
	)

	flag.Usage = func() {
		bold := func(s string) string {
			return "\033[1m" + s + "\033[0m"
		}

		fmt.Fprintf(os.Stderr, `Usage:
magic-log [flags]
magic-log config [get|set|unset] <key> [value]

%s
`, bold("Flags:"))

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, `
%s
  The CLI reads config from ~/.magiclogrc by default.
  You can override the config path using the MAGIC_LOG_CONFIG environment variable.

%s
  MAGIC_LOG_CONFIG=/path/to/custom.toml magic-log --port 4000
  magic-log config set port 4000
`, bold("Config:"), bold("Examples:"))
	}

	flag.Parse()

	if *logFormat != "json" && *logFormat != "text" && *logFormat != "csv" {
		return nil, nil, fmt.Errorf("log_format must be 'json', 'text', or 'csv'")
	}

	if *logFormat == "csv" && !*hasCSVHeader && *csvFieldsStr == "" {
		return nil, nil, fmt.Errorf("CSV format must have a header or fields passed")
	}

	cfgFile, err := Load()
	if err != nil {
		return nil, nil, err
	}

	// Track explicitly passed flags
	flagPassed := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		flagPassed[f.Name] = true
	})

	// Resolve final config
	final := &FinalConfig{
		DBFile:       resolveDBFile(*dbFile, cfgFile.Defaults.DBFile, flagPassed["db-file"]),
		Port:         pickInt(*port, cfgFile.Defaults.Port, flagPassed["port"]),
		Launch:       resolveLaunch(*launch, cfgFile.Defaults.Launch, flagPassed["launch"]),
		Echo:         *echo,
		LogFormat:    pickStr(*logFormat, cfgFile.Defaults.LogFormat, flagPassed["log-format"]),
		RegexPreset:  pickStr(*regexPreset, cfgFile.Defaults.RegexPreset, flagPassed["regex-preset"]),
		Regex:        pickStr(*parseRegex, cfgFile.Defaults.Regex, flagPassed["regex"]),
		JqFilter:     pickStr(*jqFilter, cfgFile.Defaults.JqFilter, flagPassed["jq"]),
		JqPreset:     pickStr(*jqPreset, cfgFile.Defaults.JqPreset, flagPassed["jq-preset"]),
		CSVFieldsStr: pickStr(*csvFieldsStr, cfgFile.Defaults.CSVFields, flagPassed["csv-fields"]),
		HasCSVHeader: *hasCSVHeader || cfgFile.Defaults.HasCSVHeader,
		AutoAnalyze:  !*noAutoAnalyze,
		ShowVersion:  *showVersion,
		ListPresets:  *listPresets,
	}

	return final, cfgFile, nil
}

func pickStr(cli string, def string, passed bool) string {
	if passed || cli != "" {
		return cli
	}
	return def
}

func pickInt(cli int, def int, passed bool) int {
	if passed {
		return cli
	}
	if def != 0 {
		return def
	}
	return cli
}

func pickBool(cli bool, def bool, passed bool) bool {
	if passed {
		return cli
	}
	return def
}

func resolveLaunch(cli bool, def bool, passed bool) bool {
	if passed {
		return cli
	}
	return def
}

func resolveDBFile(cli string, def string, passed bool) string {
	if passed {
		return cli
	}
	if def != "" {
		return def
	}
	return cli
}
