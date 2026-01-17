package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type zkp struct {
	r1cs constraint.ConstraintSystem
	pk   groth16.ProvingKey
	vk   groth16.VerifyingKey
}

func newZkp(c frontend.Circuit) (zkp, error) {
	// 1. Compile the circuit (R1CS Transformation)
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, c)
	if err != nil {
		return zkp{}, fmt.Errorf("frontend.Compile: %v", err)
	}

	// 2. Groth16 setup
	// Generate the Proving Key (pk) and the Verifying Key (vk)
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return zkp{}, fmt.Errorf("groth16.Setup: %v", err)
	}
	return zkp{
		r1cs: r1cs,
		pk:   pk,
		vk:   vk,
	}, nil
}

func (zkp zkp) Proof(c frontend.Circuit) groth16.Proof {
	//"Full Witness" : Compule all intermediates variables of the circuit
	witness, err := frontend.NewWitness(c, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	// 3. Generate the proof with the proving key (pk)
	proof, err := groth16.Prove(zkp.r1cs, zkp.pk, witness)
	if err != nil {
		panic(err)
	}
	return proof
}

func (zkp zkp) Verify(proof groth16.Proof, publicAssignment frontend.Circuit) error {
	// "Public Witness" : Extract what is public
	publicWitness, _ := frontend.NewWitness(publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())

	// Verify the proof with the verifier key (vk)
	return groth16.Verify(proof, zkp.vk, publicWitness)
}
