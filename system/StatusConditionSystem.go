package system

import (
	"goblin-town/component"
	"goblin-town/world"
)

type StatusConditionSystem struct {
}

var statusConditions = []string{"Poisoned"}

// StatusConditionSystem .
func (s StatusConditionSystem) Update(planets map[string]*world.Planet) map[string]*world.Planet {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			for _, entity := range level.Entities {
				for _, statusCondition := range statusConditions {
					if entity.HasComponent(statusCondition + "Component") {
						pc := entity.GetComponent(statusCondition + "Component").(component.DecayingComponent)

						if pc.Decay() {
							entity.RemoveComponent(statusCondition + "Component")
						}
					}
				}
			}
		}
	}
	return planets
}
