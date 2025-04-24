package config

import (
	"flag"
)

type FinalConfig struct {
	DBFile      string
	Port        int
	Launch      bool
	LogFormat   string
	ParsePreset string
	ParseRegex  string
	ShowVersion bool
	ListPresets bool
}

func ParseArgsAndConfig() (*FinalConfig, *Config, error) {
	// Define CLI flags
	var (
		dbFile      = flag.String("db-file", "", "Path to a DuckDB database file.")
		port        = flag.Int("port", 3000, "Port to serve the web UI on.")
		launch      = flag.Bool("launch", false, "Open the UI in a browser.")
		logFormat   = flag.String("log-format", "json", "Log format: json or text.")
		parseRegex  = flag.String("parse-regex", "", "Custom regex to parse logs.")
		parsePreset = flag.String("parse-preset", "", "Regex preset to use.")
		showVersion = flag.Bool("version", false, "Print version and exit.")
		listPresets = flag.Bool("list-presets", false, "List available presets.")
	)

	flag.Parse()

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
		DBFile:      pickStr(*dbFile, cfgFile.Defaults.DBFile, flagPassed["db-file"]),
		Port:        pickInt(*port, cfgFile.Defaults.Port, flagPassed["port"]),
		Launch:      pickBool(*launch, cfgFile.Defaults.Launch, flagPassed["launch"]),
		LogFormat:   pickStr(*logFormat, cfgFile.Defaults.LogFormat, flagPassed["log-format"]),
		ParsePreset: pickStr(*parsePreset, cfgFile.Defaults.ParsePreset, flagPassed["parse-preset"]),
		ParseRegex:  pickStr(*parseRegex, cfgFile.Defaults.ParseRegex, flagPassed["parse-regex"]),
		ShowVersion: *showVersion,
		ListPresets: *listPresets,
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
