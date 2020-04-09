package system

import (
	"goblin-town/entity"
	"goblin-town/world"
)

// System base system interface
type System interface {
	Update(*world.Level, *entity.Entity) *world.Level
}
