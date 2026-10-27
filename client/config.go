package main

import (
	"bachelor/implementation/proving/backend"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
)

type Config struct {
	Mode,
	Server,
	Circuit,
	Backend,
	Curve,
	SetupName,
	CircuitFile,
	InputsFile string
}

func parseConfig(args []string) (Config, error) {
	var config Config
	fs := flag.NewFlagSet("client", flag.ContinueOnError)

	fs.StringVar(&config.Mode, "mode", "", "local or remote")
	fs.StringVar(&config.Server, "server", "", "server address")
	fs.StringVar(&config.Circuit, "circuit", "", "built-in circuit ID (local or remote)")
	fs.StringVar(&config.Backend, "backend", "", "prove system")
	fs.StringVar(&config.Curve, "curve", "", "curve")
	fs.StringVar(&config.SetupName, "setup-name", "", "saved setup name; omit to create a fresh setup")
	fs.StringVar(&config.CircuitFile, "circuit-file", "", "single .go circuit file")
	fs.StringVar(&config.InputsFile, "inputs", "", "JSON witness file")

	fs.Usage = func() { printUsage(fs) }

	if err := fs.Parse(args); err != nil {
		return config, err
	}

	if fs.NArg() != 0 {
		return config, fmt.Errorf("unexpected positional arguments")
	}

	if config.Mode != "local" && config.Mode != "remote" {
		return config, fmt.Errorf("mode must be explicit: local or remote")
	}

	if config.Backend == "" || config.Curve == "" || config.InputsFile == "" {
		return config, fmt.Errorf("backend, curve and inputs are required")
	}

	if err := backend.ValidateBackend(config.Backend, config.Curve); err != nil {
		return config, err
	}

	if config.CircuitFile != "" && strings.ToLower(filepath.Ext(config.CircuitFile)) != ".go" {
		return config, fmt.Errorf("custom circuit must be a single .go file")
	}

	if (config.Circuit == "") == (config.CircuitFile == "") {
		return config, fmt.Errorf("specify exactly one of --circuit and --circuit-file")
	}
	if config.CircuitFile != "" && config.SetupName != "" {
		return config, fmt.Errorf("--circuit-file does not accept --setup-name")
	}

	if config.Mode == "local" {
		if config.Server != "" {
			return config, fmt.Errorf("server is not used in local modes")
		}
	} else if config.Server == "" {
		return config, fmt.Errorf("remote mode requires server address")
	}
	return config, nil
}
