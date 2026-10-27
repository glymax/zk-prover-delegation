package main

import (
	"bachelor/implementation/proving/backend"
	"flag"
	"fmt"
)

// SetupConfig contains the command-line parameters for setup preparation.
type SetupConfig struct {
	Circuit   string
	Backend   string
	Curve     string
	SetupName string
}

func parseConfig(args []string) (SetupConfig, error) {
	var setupConfig SetupConfig
	fs := flag.NewFlagSet("prepare", flag.ContinueOnError)
	fs.StringVar(&setupConfig.Circuit, "circuit", "", "fixed built-in circuit ID")
	fs.StringVar(&setupConfig.Backend, "backend", "", "backend groth16 or plonk")
	fs.StringVar(&setupConfig.Curve, "curve", "", "curve bn254, bls12-377, bls12-381 or bw6-761")
	fs.StringVar(&setupConfig.SetupName, "setup-name", "", "new setup name")
	if err := fs.Parse(args); err != nil {
		return setupConfig, err
	}
	if fs.NArg() != 0 {
		return setupConfig, fmt.Errorf("unexpected positional arguments")
	}

	if setupConfig.Circuit == "" || setupConfig.Backend == "" || setupConfig.Curve == "" || setupConfig.SetupName == "" {
		return setupConfig, fmt.Errorf("circuit, backend, curve and setup name are required")
	}
	if err := backend.ValidateBackend(setupConfig.Backend, setupConfig.Curve); err != nil {
		return setupConfig, err
	}
	return setupConfig, nil
}
