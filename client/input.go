package main

import (
	"bachelor/implementation/protocol"
	"bachelor/implementation/proving"
	"encoding/json"
	"fmt"
	"os"
)

// LoadJob reads witness JSON; custom source stays separate from the job.
func LoadJob(config Config) (proving.Job, error) {
	job := proving.Job{Circuit: config.Circuit, Backend: config.Backend, Curve: config.Curve, SetupName: config.SetupName}
	var err error
	job.Inputs, err = os.ReadFile(config.InputsFile)
	if err != nil {
		return job, err
	}
	if !json.Valid(job.Inputs) {
		return job, fmt.Errorf("inputs must contain valid JSON")
	}
	return job, nil
}

// LoadRequest is used only by the remote client.
func LoadRequest(config Config) (protocol.Request, error) {
	job, err := LoadJob(config)
	if err != nil {
		return protocol.Request{}, err
	}
	request := protocol.Request{Type: "prove", Circuit: job.Circuit, Backend: job.Backend, Curve: job.Curve, SetupName: job.SetupName, Inputs: job.Inputs}
	if config.CircuitFile != "" {
		request.Type = "custom_prove"
		request.Source, err = os.ReadFile(config.CircuitFile)
	}
	return request, err
}
