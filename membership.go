package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// MemberShipCircuit defines the constraint, SecretCode is one of PublicHashes
type MemberShipCircuit struct {
	// Le code secret que seul le prouveur connaît
	SecretCode frontend.Variable `gnark:",private"`

	// La liste publique des hashes autorisés
	PublicHashes [4]frontend.Variable `gnark:",public"`
}

func (c *MemberShipCircuit) Define(api frontend.API) error {
	// 1. Hacher le code secret
	mimc, _ := mimc.NewMiMC(api)
	mimc.Write(c.SecretCode)
	hashOfSecret := mimc.Sum()

	// 2. Vérifier que hashOfSecret est égal à l'un des PublicHashes
	// Indice : (h - h1) * (h - h2) * (h - h3) * (h - h4) == 0
	// Si l'un des termes est nul, le produit total est nul.
	var vars []frontend.Variable
	for i := range len(c.PublicHashes) {
		vars = append(vars, api.Sub(hashOfSecret, c.PublicHashes[i]))
	}
	api.AssertIsEqual(0, api.Mul(vars[0], vars[1], vars[2:]...))

	return nil
}
