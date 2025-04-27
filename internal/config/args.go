package config

import (
	"flag"
	"fmt"
	"os"
)

type FinalConfig struct {
	DBFile      string
	Port        int
	Launch      bool
	Echo        bool
	LogFormat   string
	ParsePreset string
	ParseRegex  string
	JqFilter    string
	AutoAnalyze bool
	ShowVersion bool
	ListPresets bool
}

func ParseArgsAndConfig() (*FinalConfig, *Config, error) {
	var (
		dbFile        = flag.String("db-file", "", "Path to a DuckDB database file.")
		noDBFile      = flag.Bool("no-db-file", false, "Force in-memory DB even if config has db_file.")
		port          = flag.Int("port", 3000, "Port to serve the web UI on.")
		launch        = flag.Bool("launch", false, "Open the UI in a browser.")
		noLaunch      = flag.Bool("no-launch", false, "Disable UI auto-launch (overrides config).")
		echo          = flag.Bool("echo", false, "Echo parsed stdin input to stdout")
		noAutoAnalyze = flag.Bool("no-auto-analyze", false, "Disable automatic ANALYZE of logs table")
		logFormat     = flag.String("log-format", "json", "Log format: json or text.")
		parseRegex    = flag.String("parse-regex", "", "Custom regex to parse logs.")
		parsePreset   = flag.String("parse-preset", "", "Regex preset to use.")
		jqFilter      = flag.String("jq-filter", "", "A jq expression to apply to parsed logs")
		listPresets   = flag.Bool("list-presets", false, "List available presets and exit.")
		showVersion   = flag.Bool("version", false, "Print version and exit.")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
magic-log [flags]
magic-log config [get|set|unset] <key> [value]

Flags:
`)
		flag.PrintDefaults()
	}

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
		DBFile:      resolveDBFile(*dbFile, *noDBFile, cfgFile.Defaults.DBFile, flagPassed["db-file"]),
		Port:        pickInt(*port, cfgFile.Defaults.Port, flagPassed["port"]),
		Launch:      resolveLaunch(*launch, *noLaunch, cfgFile.Defaults.Launch, flagPassed["launch"]),
		Echo:        *echo,
		LogFormat:   pickStr(*logFormat, cfgFile.Defaults.LogFormat, flagPassed["log-format"]),
		ParsePreset: pickStr(*parsePreset, cfgFile.Defaults.ParsePreset, flagPassed["parse-preset"]),
		ParseRegex:  pickStr(*parseRegex, cfgFile.Defaults.ParseRegex, flagPassed["parse-regex"]),
		JqFilter:    pickStr(*jqFilter, cfgFile.Defaults.JqFilter, flagPassed["jq-filter"]),
		AutoAnalyze: !*noAutoAnalyze,
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

func resolveLaunch(cli bool, noCli bool, def bool, passed bool) bool {
	if passed {
		return cli
	}
	if noCli {
		return false
	}
	return def
}

func resolveDBFile(cli string, disable bool, def string, passed bool) string {
	if passed {
		return cli
	}
	if disable {
		return ""
	}
	if def != "" {
		return def
	}
	return cli
}
