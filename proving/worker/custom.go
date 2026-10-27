package worker

import (
	"bachelor/implementation/proving"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// RunCustom builds and runs one custom circuit, then removes custom-worker/current.
// Source contains the complete Go file. Jobs run sequentially.
func RunCustom(job proving.Job, source []byte) (proving.Result, error) {
	workspace, err := filepath.Abs("./custom-worker")
	if err != nil {
		return proving.Result{}, err
	}
	current := filepath.Join(workspace, "current")

	if err = os.RemoveAll(current); err != nil {
		return proving.Result{}, err
	}
	if err = os.MkdirAll(filepath.Join(current, "circuit"), 0700); err != nil {
		return proving.Result{}, err
	}
	defer func() {
		if cleanupErr := os.RemoveAll(current); cleanupErr != nil {
			log.Printf("cleanup current: %v", cleanupErr)
		}
	}()
	sourceFile := filepath.Join(current, "circuit", "circuit.go")
	if err = os.WriteFile(sourceFile, source, 0600); err != nil {
		return proving.Result{}, err
	}
	worker := filepath.Join(current, "worker")
	command := exec.Command("go", "build", "-mod=readonly", "-o", worker, ".")
	command.Dir = workspace
	command.Env = append(os.Environ(), "GOWORK=off")
	command.Stdout = os.Stderr
	command.Stderr = os.Stderr
	if err = command.Run(); err != nil {
		return proving.Result{}, fmt.Errorf("custom circuit build failed: %w", err)
	}
	return runProcess(worker, job)
}
