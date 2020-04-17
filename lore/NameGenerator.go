package lore

import "math/rand"

func RandomGoblinName() (name string) {
	firstNames := []string{"Garl", "Gup", "Herr", "Mur", "Fat"}
	lastNames := []string{"Dur", "Twig", "Smellytoes", "Crook'ed"}

	name = firstNames[getRandom(0, len(firstNames))] + " " + lastNames[getRandom(0, len(lastNames))]
	return
}

func getRandom(low int, high int) int {
	return (rand.Intn((high - low))) + low
}
