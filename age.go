package main

import "github.com/consensys/gnark/frontend"

// AgeCircuit defines the constraints : Age>=Limit
type AgeCircuit struct {
	// Age is private
	Age frontend.Variable `gnark:",private"`
	// Limit is public
	Limit frontend.Variable `gnark:",public"`
}

// Define the circuit's constraints
func (circuit *AgeCircuit) Define(api frontend.API) error {
	// Contraint : Age must be greater than Limit
	api.AssertIsDifferent(-1, api.Cmp(circuit.Age, circuit.Limit))
	return nil
}
