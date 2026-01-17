package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
)

// AgeCircuit defines the constraints : Age>=Limit
type AgeCircuit struct {
	// Age is private
	Age frontend.Variable `gnark:",private"`
	// Limit is public
	Limit frontend.Variable `gnark:",public"`
}

// Define the circuit's constraints
func (circuit *AgeCircuit) Define(api frontend.API) error {
	// On veut vérifier Age >= Limit
	// Ce qui revient à vérifier que (Age - Limit) est positif

	// Dans un corps fini (Field), un nombre est "positif" s'il est
	// "petit" par rapport à la taille du Field.
	// On définit souvent une limite de bits (ex: 32 bits pour l'âge).

	diff := api.Sub(circuit.Age, circuit.Limit)

	// Cette fonction force la différence à être représentable sur 32 bits.
	// Si Age < Limit, diff sera un très grand nombre (proche de la taille du Field)
	// et cette contrainte échouera.
	bits.ToBinary(api, diff, bits.WithNbDigits(32))
	return nil
}
