package main

import (
	"fmt"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// Circuit defines the constraints : a * b = c
type Circuit struct {
	// A and B are secrets (only prover knows them)
	A frontend.Variable `gnark:",private"`
	B frontend.Variable `gnark:",private"`
	// C is public (verifier knows it)
	C frontend.Variable `gnark:",public"`
}

// Define the circuit's constraints
func (circuit *Circuit) Define(api frontend.API) error {
	// Contraint : A * B must be equal to C
	api.AssertIsEqual(api.Mul(circuit.A, circuit.B), circuit.C)
	return nil
}

func main() {
	// 1. Compile the circuit (R1CS Transformation)
	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})

	// 2. Groth16 setup
	// Generate the Proving Key (pk) and the Verifying Key (vk)
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	// Optional : Save keys to now recompute them
	// pk.WriteTo(filePK)
	// vk.WriteTo(fileVK)

	// -- prover

	// 1. Create the witness with real values
	assignment := &Circuit{
		A: 3,  // Secret
		B: 5,  // Secret
		C: 15, // Public
	}

	// 2. "Full Witness" : Compule all intermediates variables of the circuit
	witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())

	// 3. Generate the proof with the proving key (pk)
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		// Error if constraints are not fullfiled (ex: 3 * 5 != 14)
		panic(err)
	}

	// -- verifier

	// 1. Create the public witness (only public field)
	publicAssignment := &Circuit{
		C: 15,
	}

	// 2. "Public Witness" : Extract what is public
	publicWitness, _ := frontend.NewWitness(publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())

	// 3. Verify the proof with the verifier key (vk)
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println("Invalid proof")
	} else {
		fmt.Println("Accepted proof : the prover knows the factors of 15")
	}
}
