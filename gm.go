package main

import (
	"fmt"
	"goblin-town/entity"
	"goblin-town/system"
	"goblin-town/world"
	"math/rand"
)

const FOOD_MINIMUM = 40
const HOSTILE_MINIMUM = 40
const FOOD_INITIAL = 200
const HOSTILE_INITIAL = 200

type GameMaster struct {
	level *world.Level
}

func (gm *GameMaster) Init(level *world.Level) {
	gm.level = level
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
	for i := 0; i < FOOD_INITIAL; i++ {
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
	for i := 0; i < HOSTILE_INITIAL; i++ {
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
}

func (gm *GameMaster) Update() {
	foodCount := 0
	hostileCount := 0
	dragonPresent := false
	//Gather stats
	for _, e := range gm.level.Entities {
		if e.HasComponent("FoodComponent") {
			foodCount++
		}

		if e.HasComponent("HostileAIComponent") {
			hostileCount++
		}

		if e.Blueprint == "dragon" {
			dragonPresent = true
		}
	}

	if foodCount < FOOD_MINIMUM {
		fmt.Println("Below minimum number of food entities... Placing food")
		//Random food
		for i := 0; i < FOOD_MINIMUM-foodCount; i++ {
			x := rand.Intn(WIDTH)
			y := rand.Intn(HEIGHT)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(WIDTH)
				y = rand.Intn(HEIGHT)
				tile = gm.level.GetTileAt(x, y)
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
				gm.level.AddEntity(food)
			}
		}
	}

	if hostileCount < HOSTILE_INITIAL {
		fmt.Println("Below minimum number of hostile entities... Placing hostiles")
		//Random snakes
		for i := 0; i < HOSTILE_INITIAL-hostileCount; i++ {
			x := rand.Intn(WIDTH)
			y := rand.Intn(HEIGHT)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(WIDTH)
				y = rand.Intn(HEIGHT)
				tile = gm.level.GetTileAt(x, y)
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
				gm.level.AddEntity(food)
			}
		}
	}

	if dragonPresent == false {
		if getRandom(0, 1000) == 0 {
			fmt.Println("A dragon has flown in!")
			x := rand.Intn(WIDTH)
			y := rand.Intn(HEIGHT)
			tile := gm.level.GetTileAt(x, y)
			tries := 0
			for tile.Type == 2 || tile.Type == 4 || gm.level.GetEntityAt(x, y) != nil {
				x = rand.Intn(WIDTH)
				y = rand.Intn(HEIGHT)
				tile = gm.level.GetTileAt(x, y)
				tries++
				if tries > 10 {
					break
				}
			}
			if tries < 10 {
				dragon, err := entity.Create("dragon", x, y)
				if err == nil {
					gm.level.AddEntity(dragon)
				}
			}
		}
	}
}

func getRandom(low int, high int) int {
	if low == high {
		return low
	}
	return (rand.Intn((high - low))) + low
}
