package main

import (
	"flag"
	"fmt"
)

func printUsage(fs *flag.FlagSet) {
	fmt.Fprint(fs.Output(), `Usage: client [--server HOST:PORT] ...   (-h or --help)
Run from the implementation directory.
Specify exactly one of --circuit NAME and --circuit-file PATH.
Without --server, proving runs locally. --circuit-file requires a fresh setup.
Input and source file paths refer to the client machine; remote circuit and setup names refer to the server.

1. Local, built-in circuit
   client --circuit ID --backend groth16 --curve bn254 [--setup-name NAME] --inputs FILE
   --circuit ID       Required built-in circuit ID.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness.
   --setup-name NAME  Optional: load an existing setup; omit for cold.

2. Local, custom circuit
   client --circuit-file SOURCE --backend groth16 --curve bn254 --inputs FILE
   --circuit-file SOURCE      Required single .go circuit file.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness.
   Custom circuits always use a fresh setup.

3. Remote, built-in circuit
   client --server HOST:PORT --circuit ID --backend groth16 --curve bn254 [--setup-name NAME] --inputs FILE
   --server HOST:PORT Required server address.
   --circuit ID       Required built-in circuit ID.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness on client.
   --setup-name NAME  Optional: load a server-side setup; omit for cold.

4. Remote, custom circuit
   client --server HOST:PORT --circuit-file SOURCE --backend groth16 --curve bn254 --inputs FILE
   --server HOST:PORT Required server address.
   --circuit-file SOURCE      Required single .go circuit file.
   --backend NAME     Required; groth16 or plonk.
   --curve NAME       Required; bn254, bls12-377, bls12-381 or bw6-761.
   --inputs FILE      Required JSON witness on client.
   Custom circuits always use a fresh setup.
`)
}
