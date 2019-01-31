package system

import (
	"fmt"
	"goblin-town/component"
	"goblin-town/world"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var Tile_Size_W = 16
var Tile_Size_H = 16

const Sprite_Size_H = 16
const Sprite_Size_W = 16
const Window_W = 800
const Window_H = 600

type entityView struct {
	X, Y             int
	SpriteX, SpriteY int32
	Dir              int
}

var renderer *sdl.Renderer
var characterTexture *sdl.Texture
var worldTexture *sdl.Texture
var uiTexture *sdl.Texture
var window *sdl.Window

var CameraX = 0
var CameraY = 0

var Zoom = 1

func RenderSystemInit() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("Tiles", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		Window_W, Window_H, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return
	}

	image, err := img.Load("tiny_dungeon_monsters.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		return
	}

	characterTexture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return
	}

	image.Free()

	image, err = img.Load("tiny_dungeon_world.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		return
	}

	worldTexture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return
	}

	image.Free()

	image, err = img.Load("tiny_dungeon_interface.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		return
	}

	uiTexture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return
	}
}

func RenderSystemCleanup() {
	sdl.Quit()
	window.Destroy()
	renderer.Destroy()
	worldTexture.Destroy()
	characterTexture.Destroy()
	uiTexture.Destroy()

}

// RenderSystem .
func RenderSystem(planets map[string]*world.Planet) {
	if Keyboard.GetKey("a") == 1 && CameraX > 0 {
		CameraX--
	}
	if Keyboard.GetKey("d") == 1 && CameraX < Window_W {
		CameraX++
	}
	if Keyboard.GetKey("w") == 1 && CameraY > 0 {
		CameraY--
	}
	if Keyboard.GetKey("s") == 1 && CameraY < Window_H {
		CameraY++
	}

	if Keyboard.GetKey("1") == 1 {
		Tile_Size_H = 32
		Tile_Size_W = 32
	}
	if Keyboard.GetKey("2") == 1 {
		Tile_Size_H = 16
		Tile_Size_W = 16
	}
	if Keyboard.GetKey("3") == 1 {
		Tile_Size_H = 8
		Tile_Size_W = 8
	}

	if Keyboard.GetKey("4") == 1 {
		Tile_Size_H = 4
		Tile_Size_W = 4
	}

	level := planets["hub"].Levels[0]
	var seeableEntities []entityView
	for _, entity := range level.Entities {
		if entity.HasComponent("AppearanceComponent") {
			ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			dir := 0
			if entity.HasComponent("DirectionComponent") {
				dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
				dir = dc.Direction
			}
			ev := entityView{X: pc.X, Y: pc.Y, SpriteX: ac.SpriteX, SpriteY: ac.SpriteY, Dir: dir}
			seeableEntities = append(seeableEntities, ev)
		}
	}

	viewWidth := Window_W / Tile_Size_W
	viewHeight := Window_H / Tile_Size_H

	pX := Mouse.X/Tile_Size_W + CameraX
	pY := Mouse.Y/Tile_Size_H + CameraY

	view := level.GetView(CameraX, CameraY, viewWidth, viewHeight, false, false)

	renderer.Clear()
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := int32(x * Tile_Size_W)
			tY := int32(y * Tile_Size_H)
			tile := view[x][y]

			drawSprite(tX, tY, tile.SpriteX, tile.SpriteY, worldTexture) //Tile itself

			if tile == nil {
				drawSprite(tX, tY, 0, 112, worldTexture) //Empty space
			} else {
				for _, entity := range seeableEntities {
					if entity.X == tile.X && entity.Y == tile.Y {
						drawSprite(tX, tY, entity.SpriteX+(int32(entity.Dir)*Sprite_Size_W), entity.SpriteY, characterTexture) //Entity
					}
				}

				if pX == tile.X && pY == tile.Y {
					drawSprite(tX, tY, 128, 128, uiTexture) //Cursor?
				}
			}
		}
	}

	renderer.Present()
	sdl.Delay(16)
}

func drawSprite(x int32, y int32, sx int32, sy int32, texture *sdl.Texture) {
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	renderer.Copy(texture, &src, &dst)
}
