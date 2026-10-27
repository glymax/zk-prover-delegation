package main

import (
	"bachelor/implementation/proving/worker"
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark/logger"
)

func main() {
	logger.SetOutput(os.Stderr)
	var err error
	if len(os.Args) == 1 {
		err = worker.Stdio(nil)
	} else if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("prover: read one Job from stdin and write one Result to stdout.")
	} else {
		err = fmt.Errorf("expected no arguments (stdin job); use bin/prepare to create saved setups")
	}
	if err != nil {
		log.Fatal(err)
	}
}
