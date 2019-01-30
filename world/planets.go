package world

type Planet struct {
	Levels []*Level
}

func NewPlanet() (planet *Planet) {
	planet = &Planet{Levels: []*Level{}}
	planet.Levels = append(planet.Levels, NewOverworldSection(1000, 1000))
	return
}
