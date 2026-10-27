package main

import (
	custom "bachelor/implementation/custom-worker/current/circuit"
	"bachelor/implementation/proving/worker"
	"log"
	"os"

	"github.com/consensys/gnark/logger"
)

func main() {
	logger.SetOutput(os.Stderr)
	if err := worker.Stdio(custom.NewDefinition); err != nil {
		log.Fatal(err)
	}
}
