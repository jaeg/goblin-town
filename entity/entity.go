package entity

import (
	"github.com/jaeg/goblin-town/component"
)

// Entity .
type Entity struct {
	Components      map[string]component.Component
	InteractingWith component.Component
	Shown           bool
	Blueprint       string
}

func (entity *Entity) AddComponent(c component.Component) {
	if entity.Components == nil {
		entity.Components = make(map[string]component.Component)
	}

	entity.Components[c.GetType()] = c
}

func (entity *Entity) HasComponent(name string) bool {
	if entity.Components == nil {
		entity.Components = make(map[string]component.Component)
	}

	return entity.Components[name] != nil
}

func (entity *Entity) GetComponent(name string) component.Component {
	if entity.Components == nil {
		entity.Components = make(map[string]component.Component)
	}

	return entity.Components[name]
}

func (entity *Entity) RemoveComponent(name string) {
	if entity.Components == nil {
		entity.Components = make(map[string]component.Component)
	}

	entity.Components[name] = nil
}
