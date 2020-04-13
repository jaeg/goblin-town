package system

import (
	"github.com/jaeg/goblin-town/world"

	"github.com/jaeg/goblin-town/entity"
)

// System base system interface
type System interface {
	Update(*world.Level, *entity.Entity) *world.Level
}
