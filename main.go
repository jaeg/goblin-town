package main

import (
	"log"
	"math/rand"
	"time"

	"goblin-town/component"
	"goblin-town/entity"
	"goblin-town/system"
	"goblin-town/world"
)

var planets map[string]*world.Planet

func main() {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	planets = make(map[string]*world.Planet)
	planets["hub"] = world.NewPlanet()
	elapsed := time.Since(start)
	log.Printf("Generating the world took %s", elapsed)

	for i := 0; i < 10; i++ {
		entity := entity.Entity{}
		x := 1
		y := 1
		if i != 0 {
			x = rand.Intn(30)
			y = rand.Intn(30)
			message := []string{"Hello there!", "Like my hat?", "It's dangerous out here at night."}
			entity.AddComponent(&component.InteractComponent{Message: message})
			entity.AddComponent(&component.AppearanceComponent{SpriteX: 0, SpriteY: 112})
		} else {
			entity.AddComponent(&component.ShopComponent{ItemsForSale: []string{"Sword", "Bow", "Shield", "Meat"}})
			entity.AddComponent(&component.AppearanceComponent{SpriteX: 64, SpriteY: 0})
		}

		entity.AddComponent(&component.WanderAIComponent{})
		entity.AddComponent(&component.InitiativeComponent{DefaultValue: 32, Ticks: 1})
		entity.AddComponent(&component.PositionComponent{X: x, Y: y, Level: 0})
		entity.AddComponent(&component.DirectionComponent{Direction: 0})
		entity.AddComponent(&component.SolidComponent{})

		//entities = append(entities, &newPlayerEntity)
		planets["hub"].Levels[0].AddEntity(&entity)

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
