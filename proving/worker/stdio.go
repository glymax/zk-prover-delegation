package worker

import (
	"bachelor/implementation/circuits"
	"bachelor/implementation/codec"
	"bachelor/implementation/proving"
	"fmt"
	"io"
	"os"
)

// Stdio reads one Job and writes one Result. nil selects the built-in registry;
// a constructor selects the circuit compiled into a custom worker.
func Stdio(constructor circuits.Constructor) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var job proving.Job
	if err = codec.Decode(input, &job); err != nil {
		return err
	}
	if constructor == nil {
		var ok bool
		constructor, ok = circuits.Lookup(job.Circuit)
		if !ok {
			return fmt.Errorf("unknown built-in circuit %q", job.Circuit)
		}
	}
	result, err := proving.Prove(job, constructor())
	if err != nil {
		return err
	}
	output, err := codec.Encode(result)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(output)
	return err
}
