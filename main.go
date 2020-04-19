package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"goblin-town/entity"
	"goblin-town/system"
	"goblin-town/world"
)

var planets map[string]*world.Planet

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

	//Starting population of goblins
	x := rand.Intn(WIDTH)
	y := rand.Intn(HEIGHT)
	tile := level.GetTileAt(x, y)

	for tile.Type == 2 || tile.Type == 4 || level.GetEntityAt(x, y) != nil {
		x = rand.Intn(WIDTH)
		y = rand.Intn(HEIGHT)
		tile = level.GetTileAt(x, y)
	}
	level.CreateClusterOfGoblins(x, y, STARTING_GOBLINS)
	system.CenterCamera(x, y, level)

	fmt.Println("Placing food")
	//Random food
	for i := 0; i < 200; i++ {
		x := rand.Intn(WIDTH)
		y := rand.Intn(HEIGHT)
		tile := level.GetTileAt(x, y)
		tries := 0
		for tile.Type == 2 || tile.Type == 4 || level.GetEntityAt(x, y) != nil {
			x = rand.Intn(WIDTH)
			y = rand.Intn(HEIGHT)
			tile = level.GetTileAt(x, y)
			tries++
			if tries > 10 {
				break
			}
		}
		if tries > 10 {
			continue
		}
		food, err := entity.Create("rat", x, y)
		if err == nil {
			level.AddEntity(food)
		}
	}

	fmt.Println("Placing snakes")
	//Random snakes
	for i := 0; i < 200; i++ {
		x := rand.Intn(WIDTH)
		y := rand.Intn(HEIGHT)
		tile := level.GetTileAt(x, y)
		tries := 0
		for tile.Type == 2 || tile.Type == 4 || level.GetEntityAt(x, y) != nil {
			x = rand.Intn(WIDTH)
			y = rand.Intn(HEIGHT)
			tile = level.GetTileAt(x, y)
			tries++
			if tries > 10 {
				break
			}
		}
		if tries > 10 {
			continue
		}
		food, err := entity.Create("snake", x, y)
		if err == nil {
			level.AddEntity(food)
		}
	}

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
	rs := system.RenderSystem{}
	rs.Init()
	defer rs.Cleanup()
	//systems = append(systems, rs)

	//StatusCondition System
	cs := system.CleanUpSystem{}

	// Input system
	system.InputSystemInit()
	ticker := time.NewTicker(time.Second / 32)
	ticks := 0
	for _ = range ticker.C {
		ticks++
		if ticks >= 30 {
			level.NextHour()
			fmt.Println("The hour is now:", level.Hour)
			ticks = 0
		}
		//start := time.Now()
		system.InputSystem()

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
