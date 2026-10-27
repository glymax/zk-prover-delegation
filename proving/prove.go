package proving

import (
	"bachelor/implementation/proving/backend"
	"bachelor/implementation/proving/store"
	"bytes"
	"fmt"
	"github.com/consensys/gnark/frontend"
)

func Prove(job Job, definition frontend.Circuit) (Result, error) {
	if definition == nil {
		return Result{}, fmt.Errorf("constructor returned nil circuit")
	}
	full, err := decodeWitness(definition, job.Inputs, job.Curve)
	if err != nil {
		return Result{}, fmt.Errorf("witness: %w", err)
	}
	var setup *backend.Prepared
	if job.SetupName == "" {
		setup, err = backend.BuildSetup(definition, job.Backend, job.Curve)
	} else {
		path := store.Path(job.Circuit, job.Curve, job.Backend, job.SetupName)
		setup, err = store.Load(path, job.Backend, job.Curve)
	}
	if err != nil {
		return Result{}, err
	}
	proof, err := backend.Prove(setup, full)
	if err != nil {
		return Result{}, fmt.Errorf("prove: %w", err)
	}
	public, err := full.Public()
	if err != nil {
		return Result{}, err
	}
	var pb, wb, kb bytes.Buffer
	if _, err = proof.WriteTo(&pb); err != nil {
		return Result{}, err
	}
	if _, err = public.WriteTo(&wb); err != nil {
		return Result{}, err
	}
	if _, err = setup.VerifyingKey.WriteTo(&kb); err != nil {
		return Result{}, err
	}
	return Result{Circuit: job.Circuit, Backend: job.Backend, Curve: job.Curve, Proof: pb.Bytes(), PublicWitness: wb.Bytes(), VerifyingKey: kb.Bytes()}, nil
}
