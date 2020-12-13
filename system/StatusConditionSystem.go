package system

import (
	"github.com/jaeg/goblin-town/component"
	"github.com/jaeg/goblin-town/entity"
	"github.com/jaeg/goblin-town/world"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned", "Alerted"}

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
