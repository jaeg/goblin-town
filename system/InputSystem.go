package system

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type mouse struct {
	X, Y   int
	Button int
	State  int
}

type keyboard struct {
	Keys map[string]int
}

var Keyboard keyboard

var Mouse mouse

func InputSystemInit() {
	Keyboard.Keys = make(map[string]int, 0)
}

func InputSystem() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			os.Exit(0)
			break
		case *sdl.MouseMotionEvent:
			Mouse.X = int(t.X)
			Mouse.Y = int(t.Y)
		case *sdl.MouseButtonEvent:
			Mouse.X = int(t.X)
			Mouse.Y = int(t.Y)
			Mouse.Button = int(t.Button)
			Mouse.State = int(t.State)
		case *sdl.KeyboardEvent:
			Keyboard.Keys[string(t.Keysym.Sym)] = int(t.State)
			fmt.Println(Keyboard)
		}
	}
}
