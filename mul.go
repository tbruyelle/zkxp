package main

import "github.com/consensys/gnark/frontend"

// MulCircuit defines the constraints : a * b = c
type MulCircuit struct {
	// A and B are secrets (only prover knows them)
	A frontend.Variable `gnark:",private"`
	B frontend.Variable `gnark:",private"`
	// C is public (verifier knows it)
	C frontend.Variable `gnark:",public"`
}

// Define the circuit's constraints
func (circuit *MulCircuit) Define(api frontend.API) error {
	// Contraint : A * B must be equal to C
	api.AssertIsEqual(api.Mul(circuit.A, circuit.B), circuit.C)
	return nil
}
