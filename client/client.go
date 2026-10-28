package main

import (
	"bachelor/implementation/proving"
	"bachelor/implementation/proving/worker"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	if err = run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg Config) error {
	var result proving.Result
	mode := "local"
	if cfg.Server == "" {
		job, err := LoadJob(cfg)
		if err != nil {
			return err
		}
		if cfg.CircuitFile != "" {
			var source []byte
			source, err = os.ReadFile(cfg.CircuitFile)
			if err != nil {
				return err
			}
			result, err = worker.RunCustom(job, source)
		} else {
			result, err = worker.RunBuiltin(job)
		}
		if err != nil {
			return err
		}
	} else {
		mode = "remote"
		response, err := exchange(cfg)
		if err != nil {
			return err
		}
		result = proving.Result{Circuit: response.Circuit, Backend: response.Backend, Curve: response.Curve, Proof: response.Proof, PublicWitness: response.PublicWitness, VerifyingKey: response.VerifyingKey}
	}
	// Result stays in memory; forwarding to a verifier is the next stage.
	fmt.Printf("Result (%s): proof", mode)
	if result.Circuit != "" {
		fmt.Printf("; circuit=%s", result.Circuit)
	}
	fmt.Printf("; proof=%d bytes; public=%d bytes; vk=%d bytes\n", len(result.Proof), len(result.PublicWitness), len(result.VerifyingKey))
	return nil
}
