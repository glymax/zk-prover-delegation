package store

import (
	"bachelor/implementation/proving/backend"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Path identifies a prepared setup relative to the project directory.
func Path(circuit, curve, backend, name string) string {
	return filepath.Join("./setups", circuit, curve, backend, name)
}

func Load(path, system, curve string) (*backend.Prepared, error) {
	s, err := backend.Empty(system, curve)
	if err != nil {
		return nil, err
	}
	for name, v := range map[string]io.ReaderFrom{"constraints.bin": s.Constraints, "proving.key": s.ProvingKey, "verifying.key": s.VerifyingKey} {
		if err = ReadBinary(filepath.Join(path, name), v); err != nil {
			return nil, fmt.Errorf("load prepared setup: %w", err)
		}
	}
	return s, nil
}

func Save(path string, setup *backend.Prepared) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("setup already exists: %s", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return err
	}
	temp, err := os.MkdirTemp(parent, ".prepare-")
	if err != nil {
		return err
	}
	defer func() {
		if cleanupErr := os.RemoveAll(temp); cleanupErr != nil {
			log.Printf("cleanup setup: %v", cleanupErr)
		}
	}()
	for name, v := range map[string]io.WriterTo{"constraints.bin": setup.Constraints, "proving.key": setup.ProvingKey, "verifying.key": setup.VerifyingKey} {
		if err = WriteBinary(filepath.Join(temp, name), v); err != nil {
			return err
		}
	}
	return os.Rename(temp, path)
}

func WriteBinary(path string, value io.WriterTo) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, writeErr := value.WriteTo(file)
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func ReadBinary(path string, value io.ReaderFrom) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, readErr := value.ReadFrom(file)
	closeErr := file.Close()
	if readErr != nil {
		return readErr
	}
	return closeErr
}
