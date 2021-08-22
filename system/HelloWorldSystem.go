package system

import (
	"log"

	"github.com/jaeg/goblin-town/component"
)

// HelloWorldSystem .
type HelloWorldSystem struct {
}

// Update .
func (HelloWorldSystem) Update(a *component.HelloWorldComponent) {
	log.Println("hello world")
}
