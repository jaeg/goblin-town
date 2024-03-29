package system

import (
	"github.com/jaeg/goblin-town/component"
	"github.com/jaeg/goblin-town/entity"
	"github.com/jaeg/goblin-town/world"
)

type InitiativeSystem struct {
}

// InitiativeSystem .
func (s InitiativeSystem) Update(level *world.Level, entity *entity.Entity) *world.Level {

	if entity.HasComponent("InitiativeComponent") {
		ic := entity.GetComponent("InitiativeComponent").(*component.InitiativeComponent)
		ic.Ticks--

		if ic.Ticks <= 0 {
			ic.Ticks = ic.DefaultValue
			if ic.OverrideValue > 0 {
				ic.Ticks = ic.OverrideValue
			}

			if !entity.HasComponent("MyTurnComponent") {

				//Handle sleep schedules.
				canGo := false
				if entity.HasComponent("NocturnalComponent") {
					if level.Hour >= 20 || (level.Hour >= 0 && level.Hour <= 7) {
						canGo = true
					}
				} else {
					if level.Hour <= 20 && level.Hour >= 7 {
						canGo = true
					}
				}

				//Can't sleep if alerted.
				if entity.HasComponent("AlertedComponent") {
					canGo = true
				}

				//Never sleeps
				if entity.HasComponent("NeverSleepComponent") {
					canGo = true
				}

				if canGo {
					mTC := &component.MyTurnComponent{}
					entity.AddComponent(mTC)
				}
			}
		}
	}

	return level
}
