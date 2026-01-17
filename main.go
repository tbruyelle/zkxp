package main

import (
	"flag"
	"fmt"

	"github.com/consensys/gnark/frontend"
)

func prove(circuit, assignment, publicAssignement frontend.Circuit, verifiedLabel string) {
	zkp, err := newZkp(circuit)
	if err != nil {
		panic(err)
	}
	// -- prover
	//  Create the proof with real values
	proof := zkp.Proof(assignment)

	// -- verifier
	// pass the public witness (only public field)
	err = zkp.Verify(proof, publicAssignement)
	if err != nil {
		panic(fmt.Errorf("Invalid proof: %w", err))
	}
	fmt.Printf("Accepted proof : %s\n", verifiedLabel)
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {

	case "age":
		prove(
			&AgeCircuit{},
			&AgeCircuit{
				Age:   24,
				Limit: 18,
			},
			&AgeCircuit{
				Limit: 18,
			},
			"the prover knows that age is greater or equal to 18",
		)

	case "mul":
		prove(
			&MulCircuit{},
			&MulCircuit{
				A: 3,  // Secret
				B: 5,  // Secret
				C: 15, // Public
			},
			&MulCircuit{
				C: 15,
			},
			"the prover knows the factors of 15",
		)
	}
}
