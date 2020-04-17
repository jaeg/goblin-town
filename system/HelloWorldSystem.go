package system

import (
	"fmt"
	"github.com/jaeg/goblin-town/component"
)

// HelloWorldSystem .
type HelloWorldSystem struct {
}

// Update .
func (HelloWorldSystem) Update(a *component.HelloWorldComponent) {
	fmt.Println("hello world")
}
