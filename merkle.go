package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type MerkleProofCircuit struct {
	Root   frontend.Variable   `gnark:",public"`
	Secret frontend.Variable   `gnark:",private"`
	Path   []frontend.Variable `gnark:",private"`
	Helper []frontend.Variable `gnark:",private"` // Indique si on est à gauche ou à droite à chaque étage
}

func (c *MerkleProofCircuit) Define(api frontend.API) error {
	// 1. Calculer le hash de départ (la feuille)
	h, _ := mimc.NewMiMC(api)
	h.Write(c.Secret)
	currentHash := h.Sum()

	// 2. Remonter l'arbre
	for i := range len(c.Path) {
		api.AssertIsBoolean(c.Helper[i])

		// Selon Helper[i], currentHash est à gauche ou à droite
		// h(A, B) ou h(B, A)
		left := api.Select(c.Helper[i], currentHash, c.Path[i])
		right := api.Select(c.Helper[i], c.Path[i], currentHash)
		h.Reset()
		h.Write(left, right)
		currentHash = h.Sum()
	}

	// 3. Vérifier qu'on tombe sur la racine
	api.AssertIsEqual(currentHash, c.Root)
	return nil
}
