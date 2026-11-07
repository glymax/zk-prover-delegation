package main

import (
	"bachelor/implementation/circuits"
	"bachelor/implementation/proving/backend"
	"bachelor/implementation/proving/store"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark/logger"
)

func main() {
	logger.SetOutput(os.Stderr)
	setupConfig, err := parseConfig(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	if err = prepareSetup(setupConfig); err != nil {
		log.Fatal(err)
	}
}

func prepareSetup(setupConfig SetupConfig) error {

	constructor, ok := circuits.Lookup(setupConfig.Circuit)
	if !ok {
		return fmt.Errorf("unknown built-in circuit %q", setupConfig.Circuit)
	}

	path := store.Path(setupConfig.Circuit, setupConfig.Curve, setupConfig.Backend, setupConfig.SetupName)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("setup already exists: %s", path)
	} else if !os.IsNotExist(err) {
		return err
	}

	circuit := constructor()
	if circuit == nil {
		return fmt.Errorf("constructor returned nil circuit")
	}

	setup, err := backend.BuildSetup(circuit, setupConfig.Backend, setupConfig.Curve)
	if err != nil {
		return err
	}

	if err = store.Save(path, setup); err != nil {
		return err
	}

	fmt.Printf("Prepared %s\n", path)
	return nil
}
