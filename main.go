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

	//Random population of goblins
	for i := 0; i < 10; i++ {
		x := rand.Intn(10)
		y := rand.Intn(10)
		goblin, err := entity.Create("goblin", x, y)
		if err == nil {
			planets["hub"].Levels[0].AddEntity(goblin)
		}
	}

	//Random food
	for i := 0; i < 20; i++ {
		x := rand.Intn(20)
		y := rand.Intn(20)
		food, err := entity.Create("rat", x, y)
		if err == nil {
			planets["hub"].Levels[0].AddEntity(food)
		}
	}

	systems := make([]system.System, 0)

	//Initiative System
	systems = append(systems, system.InitiativeSystem{})

	//AI System
	systems = append(systems, system.AISystem{})

	//StatusCondition System
	systems = append(systems, system.StatusConditionSystem{})

	//Render System
	rs := system.RenderSystem{}
	rs.Init()
	defer rs.Cleanup()
	systems = append(systems, rs)

	//StatusCondition System
	systems = append(systems, system.CleanUpSystem{})

	// Input system
	system.InputSystemInit()
	ticker := time.NewTicker(time.Second / 32)

	for _ = range ticker.C {
		//start := time.Now()
		system.InputSystem()
		for s := range systems {
			planets = systems[s].Update(planets)
		}

		//elapsed := time.Since(start)
		//log.Printf("Game loop took %s", elapsed)
	}

}
