package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jaeg/goblin-town/entity"
	"github.com/jaeg/goblin-town/render"
	"github.com/jaeg/goblin-town/system"
	"github.com/jaeg/goblin-town/world"
)

const STARTING_GOBLINS = 5
const STARTING_GOBLIN_CLUSTERS = 10

const WIDTH = 200
const HEIGHT = 200

func main() {
	entity.FactoryLoad("entities.blueprints")
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	level := world.NewOverworldSection(WIDTH, HEIGHT)
	elapsed := time.Since(start)
	log.Printf("Generating the world took %s", elapsed)

	gm := &GameMaster{}
	gm.Init(level)

	systems := make([]system.System, 0)

	//Initiative System
	systems = append(systems, system.InitiativeSystem{})

	//AI System
	systems = append(systems, system.AISystem{})

	//Goblin System
	systems = append(systems, system.GoblinAISystem{})

	//StatusCondition System
	systems = append(systems, system.StatusConditionSystem{})

	//Render System
	rs := render.Renderer{}
	rs.Init()
	defer rs.Cleanup()
	//systems = append(systems, rs)

	//StatusCondition System
	cs := system.CleanUpSystem{}

	// Input system
	system.InputSystemInit()
	ticker := time.NewTicker(time.Second / 32)
	ticks := 0
	for range ticker.C {
		ticks++
		if ticks >= 120 {
			level.NextHour()
			fmt.Println("The hour is now:", level.Hour)
			ticks = 0
		}

		if ticks%30 == 0 {
			if render.Beat == 1 {
				render.Beat = 0
			} else {
				render.Beat = 1
			}
		}
		//start := time.Now()
		system.InputSystem()

		gm.Update()

		for _, entity := range level.Entities {
			for s := range systems {
				level = systems[s].Update(level, entity)
			}
		}

		level = cs.Update(level)
		level = rs.Update(level)

		//elapsed := time.Since(start)
		//log.Printf("Game loop took %s", elapsed)
	}

}
