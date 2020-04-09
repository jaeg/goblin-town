package system

import (
	"goblin-town/component"
	"goblin-town/entity"
	"goblin-town/world"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned"}

// StatusConditionSystem .
func (s StatusConditionSystem) Update(level *world.Level, entity *entity.Entity) *world.Level {

	for _, statusCondition := range statusConditions {
		if entity.HasComponent(statusCondition + "Component") {
			pc := entity.GetComponent(statusCondition + "Component").(component.DecayingComponent)

			if pc.Decay() {
				entity.RemoveComponent(statusCondition + "Component")
			}
		}
	}

	return level
}
