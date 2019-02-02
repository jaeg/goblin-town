package main

import (
	"log"
	"math/rand"
	"time"

	"goblin-town/entity"
	"goblin-town/system"
	"goblin-town/world"
)

var planets map[string]*world.Planet

func main() {
	entity.FactoryLoad("entities.blueprints")
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	planets = make(map[string]*world.Planet)
	planets["hub"] = world.NewPlanet()
	elapsed := time.Since(start)
	log.Printf("Generating the world took %s", elapsed)

	for i := 0; i < 10; i++ {
		x := rand.Intn(30)
		y := rand.Intn(30)
		goblin, err := entity.Create("goblin", x, y)
		if err == nil {
			planets["hub"].Levels[0].AddEntity(goblin)
		}
	}

	system.RenderSystemInit()
	system.InputSystemInit()
	defer system.RenderSystemCleanup()

	ticker := time.NewTicker(time.Second / 32)

	for _ = range ticker.C {
		//start := time.Now()
		system.InputSystem()
		system.InitiativeSystem(planets)
		system.AISystem(planets)
		system.RenderSystem(planets)
		system.StatusConditionSystem(planets)
		planets = system.CleanUpSystem(planets)
		//elapsed := time.Since(start)
		//log.Printf("Game loop took %s", elapsed)
	}

}
