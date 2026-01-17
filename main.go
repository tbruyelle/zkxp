package main

import (
	"flag"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
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

	case "merkle":
		// TODO

	case "membership":
		// Voici la méthode propre pour obtenir un hash compatible avec le circuit :
		hSecret := mimcHash(42)
		hOther1 := mimcHash(6789)
		hOther2 := mimcHash(9999)
		hOther3 := mimcHash(1111)

		prove(
			&MemberShipCircuit{},
			&MemberShipCircuit{
				SecretCode: 42,
				PublicHashes: [4]frontend.Variable{
					hSecret,
					hOther1,
					hOther2,
					hOther3,
				},
			},
			&MemberShipCircuit{
				PublicHashes: [4]frontend.Variable{
					hSecret,
					hOther1,
					hOther2,
					hOther3,
				},
			},
			"the provier knows the code of one of the hashes",
		)

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

func mimcHash(val int) []byte {
	f := ecc.BN254.ScalarField()
	h := hash.MIMC_BN254.New()

	// On convertit l'int en big.Int, puis on récupère ses bytes au format "Big Endian"
	// car c'est ce que le circuit attend pour un élément de Field
	var b big.Int
	b.SetInt64(int64(val))

	// On s'assure que le nombre est bien dans le Field
	res := b.Mod(&b, f)

	h.Write(res.Bytes())
	return h.Sum(nil)
}
