package system

import (
	"fmt"
	"goblin-town/component"
	"goblin-town/world"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type entityView struct {
	X, Y             int
	SpriteX, SpriteY int32
}

var renderer *sdl.Renderer
var texture *sdl.Texture
var window *sdl.Window

func RenderSystemInit() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("Tiles", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return
	}

	image, err := img.Load("goblin_cave.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		return
	}
	defer image.Free()

	texture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return
	}
}

func RenderSystemCleanup() {
	sdl.Quit()
	window.Destroy()
	renderer.Destroy()
	texture.Destroy()
}

// RenderSystem .
func RenderSystem(planets map[string]*world.Planet) {
	level := planets["hub"].Levels[0]
	var seeableEntities []entityView
	for _, entity := range level.Entities {
		if entity.HasComponent("AppearanceComponent") {
			ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			ev := entityView{X: pc.X, Y: pc.Y, SpriteX: ac.SpriteX, SpriteY: ac.SpriteY}
			seeableEntities = append(seeableEntities, ev)
		}
	}

	viewWidth := 25
	viewHeight := 19
	cX := viewWidth / 2
	cY := viewHeight / 2
	pX := Mouse.X/32 - viewWidth/2 + cX
	pY := Mouse.Y/32 - viewHeight/2 + cY
	fmt.Println(pX, pY)
	view := level.GetView(cX, cY, viewWidth, viewHeight, false)

	renderer.Clear()
	for y := 0; y < len(view[0]); y++ {
		for x := 0; x < len(view); x++ {
			tX := int32(x * 32)
			tY := int32(y * 32)
			tile := view[x][y]
			if tile == nil {
				drawSprite(tX, tY, 64, 0, texture)
			} else {
				if pX == tile.X && pY == tile.Y {
					drawSprite(tX, tY, 32, 0, texture) //Cursor?
					fmt.Println("Hi")
				} else {
					drawTile := true
					for _, entity := range seeableEntities {
						if entity.X == tile.X && entity.Y == tile.Y {
							if drawTile {
								drawSprite(tX, tY, entity.SpriteX, entity.SpriteY, texture)
								drawTile = false
							}
						}
					}
					if drawTile {
						drawSprite(tX, tY, tile.SpriteX, tile.SpriteY, texture)
					}
				}
			}
		}
	}

	renderer.Present()
	sdl.Delay(16)
}

func drawSprite(x int32, y int32, sx int32, sy int32, texture *sdl.Texture) {
	src := sdl.Rect{X: sx, Y: sy, W: 32, H: 32}
	dst := sdl.Rect{X: x, Y: y, W: 32, H: 32}
	renderer.Copy(texture, &src, &dst)
}
