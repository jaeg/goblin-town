package system

import (
	"github.com/jaeg/goblin-town/entity"
	"github.com/jaeg/goblin-town/world"
)

// System base system interface
type System interface {
	Update(*world.Level, *entity.Entity) *world.Level
}
