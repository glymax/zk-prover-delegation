package proving

// Job describes one proof generation, either cold or warm.
type Job struct {
	Circuit,
	Backend,
	Curve,
	SetupName string
	Inputs []byte
}
type Result struct {
	Circuit,
	Backend,
	Curve string
	Proof, PublicWitness, VerifyingKey []byte
}
