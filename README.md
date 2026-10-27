# zk-prover-delegation

This project is part of my bachelor's thesis at TUM on delegating zk-SNARK proving to cloud-based trusted execution environments. It is written in Go using gnark.

You can generate proofs locally or send a proving request to a server, using Groth16 or PLONK. Both modes support the included circuits and custom circuits supplied as a single Go file. Inputs are stored as JSON, and network messages are encoded with CBOR and sent over Noise NN. TEE integration is planned; the current server runs as a regular process.

## Build and run

Use Go 1.27.2 on Linux or WSL and run the commands below from the repository root. Dependency versions are specified in `go.mod`, with checksums recorded in `go.sum`.

Build the prover worker. The client and server launch this program to generate proofs for the included circuits:

```bash
go build -o bin/prover ./prover
```

Build the setup preparation tool. It creates reusable setups for runs that specify `--setup-name`:

```bash
go build -o bin/prepare ./prepare
```

Build the client. It reads the inputs and either runs proving locally or sends a request to a server:

```bash
go build -o bin/client ./client
```

Build the server, which receives remote requests and runs the prover:

```bash
go build -o bin/server ./server
```

To try a local run, generate a Groth16 proof for the Cubic circuit over BN254 using the sample inputs. Without `--setup-name`, the prover creates a fresh setup:

```bash
./bin/client --mode local --circuit cubic --backend groth16 --curve bn254 \
  --inputs client/inputs/cubic/000001.json
```

To try remote proving on the same machine, start the server:

```bash
./bin/server --listen 127.0.0.1:9000
```

Then send the same proving request from another terminal:

```bash
./bin/client --mode remote --server 127.0.0.1:9000 \
  --circuit cubic --backend groth16 --curve bn254 \
  --inputs client/inputs/cubic/000001.json
```

The server handles one connection and then exits. Use `./bin/client --help` for custom-circuit options. To reuse a setup, create it with `bin/prepare` and pass its name through `--setup-name`; otherwise a fresh setup is generated.

## Current scope

- TDX and real remote attestation are not implemented. The attestation exchange is a stub; Noise NN alone does not authenticate the server.
- PLONK uses an unsafe test SRS, intended only for experiments.
- Custom circuits provide `NewDefinition() frontend.Circuit` and implement `Define`. Inputs must match the gnark witness JSON schema. Custom code is compiled and executed without sandboxing.
- Each remote message must fit within 65,519 bytes before Noise encryption; fragmentation is not implemented.
- Proof verification and benchmark automation are planned. Sample circuits include adaptations of the [gnark examples](https://github.com/Consensys/gnark/tree/v0.16.4/examples).

`scripts/setup-node.sh` prepares an Ubuntu 22.04 amd64 experiment node; it defaults to a dry run. Build outputs, temporary workers and the separately cloned Engine framework are excluded from Git.
