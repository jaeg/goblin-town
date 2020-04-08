package system

import (
	"fmt"
	"goblin-town/world"
)

type CleanUpSystem struct {
}

// CleanUpSystem .
func (s CleanUpSystem) Update(planets map[string]*world.Planet) map[string]*world.Planet {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			for i, entity := range level.Entities {
				if entity.HasComponent("MyTurnComponent") {
					entity.RemoveComponent("MyTurnComponent")
				}

				if entity.HasComponent("DeadComponent") {
					level.Entities = append(level.Entities[:i], level.Entities[i+1:]...)
					fmt.Println("Killed")
				}

			}
		}
	}

	return planets
}
