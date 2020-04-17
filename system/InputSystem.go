package system

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type mouse struct {
	X, Y    int
	Button  int
	Clicked bool
}

type keyboard struct {
}

func (k keyboard) GetKey(key string) uint8 {
	sc := sdl.GetScancodeFromName(key)
	return sdl.GetKeyboardState()[sc]
}

var Keyboard keyboard

var Mouse mouse

func InputSystemInit() {
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
			if t.State == sdl.PRESSED {
				Mouse.Clicked = true
			} else {
				Mouse.Clicked = false
			}
		}
	}
}
