package main

import (
	"fmt"
)

func main() {
	zkp, err := newZkp(&MulCircuit{})
	if err != nil {
		panic(err)
	}
	// -- prover
	//  Create the proof with real values
	proof := zkp.Proof(&MulCircuit{
		A: 3,  // Secret
		B: 5,  // Secret
		C: 15, // Public
	})

	// -- verifier
	// pass the public witness (only public field)
	err = zkp.Verify(proof, &MulCircuit{
		C: 15,
	})
	if err != nil {
		panic(fmt.Errorf("Invalid proof: %w", err))
	}
	fmt.Println("Accepted proof : the prover knows the factors of 15")
}
