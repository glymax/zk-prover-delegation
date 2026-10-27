// Package worker builds, launches, and serves proving processes.
package worker

import (
	"bachelor/implementation/codec"
	"bachelor/implementation/proving"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// RunBuiltin launches the precompiled prover for a registered circuit.
func RunBuiltin(job proving.Job) (proving.Result, error) {
	return runProcess("./bin/prover", job)
}

func runProcess(executable string, job proving.Job) (proving.Result, error) {
	input, err := codec.Encode(job)
	if err != nil {
		return proving.Result{}, err
	}
	cmd := exec.Command(executable)
	cmd.WaitDelay = time.Second
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return proving.Result{}, fmt.Errorf("proving process failed: %w", err)
	}
	var result proving.Result
	if err = codec.Decode(output, &result); err != nil {
		return proving.Result{}, fmt.Errorf("decode result: %w", err)
	}
	return result, nil
}
