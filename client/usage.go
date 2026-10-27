package main

import (
	"flag"
	"fmt"
)

func printUsage(fs *flag.FlagSet) {
	fmt.Fprint(fs.Output(), `Usage: client --mode MODE ...   (-h or --help)
Run from the implementation directory.
Specify exactly one of --circuit NAME and --circuit-file PATH.
--server is required only in remote mode; --circuit-file requires cold setup.

1. Local, built-in circuit
   client --mode local --circuit ID --backend groth16 --curve bn254 [--setup-name NAME] --inputs FILE
   --circuit ID       Required built-in circuit ID.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness.
   --setup-name NAME  Optional: load an existing setup; omit for cold.

2. Local, custom circuit
   client --mode local --circuit-file SOURCE --backend groth16 --curve bn254 --inputs FILE
   --circuit-file SOURCE      Required single .go circuit file.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness.
    Required; source circuits always use cold.

3. Remote, built-in circuit
   client --mode remote --server HOST:PORT --circuit ID --backend groth16 --curve bn254 [--setup-name NAME] --inputs FILE
   --server HOST:PORT Required server address.
   --circuit ID       Required built-in circuit ID.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness on client.
   --setup-name NAME  Optional: load a server-side setup; omit for cold.

4. Remote, custom circuit
   client --mode remote --server HOST:PORT --circuit-file SOURCE --backend groth16 --curve bn254 --inputs FILE
   --server HOST:PORT Required server address.
   --circuit-file SOURCE      Required single .go circuit file.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness on client.
    Required; source circuits always use cold.
`)
}
