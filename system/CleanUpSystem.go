package system

import (
	"fmt"
	"github.com/jaeg/goblin-town/world"
)

type CleanUpSystem struct {
}

// CleanUpSystem .
func (s CleanUpSystem) Update(level *world.Level) *world.Level {
	for _, entity := range level.Entities {
		if entity.HasComponent("MyTurnComponent") {
			entity.RemoveComponent("MyTurnComponent")
		}

		if entity.HasComponent("DeadComponent") {
			level.RemoveEntity(entity)
			fmt.Println("Killed")
		}

	}

	return level
}
