package system

import "goblin-town/world"

// System base system interface
type System interface {
	Update(map[string]*world.Planet) map[string]*world.Planet
}
