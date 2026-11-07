# zk-prover-delegation

This project is part of my bachelor's thesis at TUM on delegating zk-SNARK proving to cloud-based trusted execution environments. It is written in Go using gnark. The [IntroTalk](docs/Introtalk.pdf) describes the proposed architecture and research plan.

You can generate proofs locally or send a proving request to a server, using Groth16 or PLONK. Both modes support built-in circuits and custom circuits supplied as a single Go file. Inputs are JSON; network messages use CBOR over Noise KK. TEE integration is planned; the current server runs as a regular process.

## Build and run

Use Go 1.27.2 on Linux or WSL and run all commands from the repository root. Dependencies are specified in `go.mod` and `go.sum`. Build each program before its first use and rebuild it after changing its source.

### Build the programs

```bash
go build -o bin/client ./client
go build -o bin/prover ./prover
go build -o bin/setupgen ./setupgen
go build -o bin/server ./server
go build -o bin/authkeygen ./authkeygen
```

| Program | Purpose |
| --- | --- |
| `client` | Read inputs and request local or remote proving. |
| `prover` | Worker launched automatically for built-in circuits. |
| `setupgen` | Create a reusable setup for a built-in circuit. |
| `server` | Handle one remote proving session. |
| `authkeygen` | Generate both participants' authentication keys. |

## Usage

### Local proving

Generate a Groth16 proof for Cubic over BN254:

```bash
./bin/client --circuit cubic --backend groth16 --curve bn254 \
  --inputs client/inputs/cubic/000001.json
```

Without `--server`, proving runs locally and requires no authentication keys. Without `--setup-name`, each run compiles the circuit and creates a fresh setup. The client reports the sizes of the resulting proof, public witness and verification key; it does not save or verify them yet.

### Reusing a setup

Create a named setup for a built-in circuit:

```bash
./bin/setupgen --circuit cubic --backend groth16 --curve bn254 \
  --setup-name example
```

This writes constraints and keys to `setups/cubic/bn254/groth16/example/`. Existing setups are not overwritten. Reuse this setup with the client built above:

```bash
./bin/client --circuit cubic --backend groth16 --curve bn254 \
  --setup-name example --inputs client/inputs/cubic/000001.json
```

The circuit, backend and curve must match the saved setup. Create a new setup after changing the circuit. For remote proving, the setup must exist on the server.

### Remote proving

Generate the client and server authentication keys:

```bash
./bin/authkeygen
```

`authkeygen` creates pinned X25519 identities to model the result of out-of-band PKI enrollment. It writes raw 32-byte keys to `client/keys` and `server/keys`: each participant receives its own `private.key` and the other participant's public key. Every invocation replaces both identities. On separate machines, distribute only the corresponding participant's key directory, and update both sides after regeneration. Programs load these fixed paths relative to the project root.

For a test on one machine, start the server:

```bash
./bin/server --listen 127.0.0.1:9000
```

From another terminal, send a request:

```bash
./bin/client --server 127.0.0.1:9000 \
  --circuit cubic --backend groth16 --curve bn254 \
  --inputs client/inputs/cubic/000001.json
```

The server launches its local `bin/prover`, returns the result and exits after one connection. Restart it for each request. On separate machines, use the server's reachable address; source and input paths are always resolved on the client.

### Custom circuits

Use the same client with `--circuit-file` instead of `--circuit`. For example, this circuit sums three inputs:

```bash
./bin/client --circuit-file client/custom-circuits/sum/circuit.go \
  --backend groth16 --curve bn254 \
  --inputs client/inputs/sum/size-3.json
```

The Go file must use an importable package (not `main`), implement `Define`, and export `NewDefinition() frontend.Circuit`. The constructor fixes the circuit's shape, including slice lengths; JSON supplies witness values and does not change that shape.

A custom worker is compiled automatically using the project's dependencies, so no separate worker build command is needed. Go must be available on the machine performing the proof. Custom circuits always create a fresh setup and do not accept `--setup-name`. Run only one custom job at a time per checkout.

For remote execution, start the server as above, then send the source and inputs with:

```bash
./bin/client --server 127.0.0.1:9000 \
  --circuit-file client/custom-circuits/sum/circuit.go \
  --backend groth16 --curve bn254 \
  --inputs client/inputs/sum/size-3.json
```

The server builds and runs the custom worker. Use `./bin/client --help` for the full list of options.

## Current scope

- TDX and real remote attestation are not implemented. The attestation exchange is a stub. Noise KK authenticates possession of pinned keys, but does not verify a TEE or bind attestation evidence to the session yet.
- PLONK uses an unsafe test SRS, intended only for experiments.
- Inputs must match the gnark witness JSON schema. Custom code is compiled and executed without sandboxing.
- Each remote message must fit within 65,519 bytes before Noise encryption; fragmentation is not implemented.
- Proof verification and benchmark automation are planned.

`scripts/setup-node.sh` prepares an Ubuntu 22.04 amd64 experiment node; it defaults to a dry run. Build outputs, temporary workers and the separately cloned Engine framework are excluded from Git.

Sample circuits include adaptations of the [gnark examples](https://github.com/Consensys/gnark/tree/v0.16.4/examples). Their Apache 2.0 license is preserved in [docs/LICENSE](docs/LICENSE); this file does not assign a license to the project's own code.
