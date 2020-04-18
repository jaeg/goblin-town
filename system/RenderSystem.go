package system

import (
	"fmt"
	"goblin-town/component"
	"goblin-town/world"
	"os"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var Tile_Size_W = 32
var Tile_Size_H = 32

const Sprite_Size_H = 16
const Sprite_Size_W = 16
const Window_W = 1000
const Window_H = 576
const World_W = 800
const World_H = 576

type RenderSystem struct {
}

var renderer *sdl.Renderer
var characterTexture *sdl.Texture
var worldTexture *sdl.Texture
var uiTexture *sdl.Texture
var window *sdl.Window
var font *ttf.Font

var CameraX = 0
var CameraY = 0

var Zoom = 1

var releasedZoom = true

func (s RenderSystem) Init() {
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

	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to initialize TTF: %s\n", err)
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

	if font, err = ttf.OpenFont("Roboto-Regular.ttf", 30); err != nil {
		fmt.Printf("Failed to open font: %s\n", err)
		return
	}
}

func (s RenderSystem) Cleanup() {
	sdl.Quit()
	window.Destroy()
	renderer.Destroy()
	worldTexture.Destroy()
	characterTexture.Destroy()
	uiTexture.Destroy()

}

// RenderSystem .
func (s RenderSystem) Update(level *world.Level) *world.Level {
	pX := Mouse.X/Tile_Size_W + CameraX
	pY := Mouse.Y/Tile_Size_H + CameraY

	if Keyboard.GetKey("a") == 1 && CameraX > 0 {
		CameraX--
	}
	if Keyboard.GetKey("d") == 1 && CameraX < level.Width-World_W/Tile_Size_W-1 {
		CameraX++
	}
	if Keyboard.GetKey("w") == 1 && CameraY > 0 {
		CameraY--
	}
	if Keyboard.GetKey("s") == 1 && CameraY < level.Height-World_H/Tile_Size_H-1 {
		CameraY++
	}

	if Keyboard.GetKey("1") == 1 {
		if releasedZoom == true {
			Tile_Size_H = 32
			Tile_Size_W = 32
			CameraX = Mouse.X/Tile_Size_W + CameraX
			CameraY = Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}

	}
	if Keyboard.GetKey("2") == 1 {
		if releasedZoom == true {
			Tile_Size_H = 16
			Tile_Size_W = 16
			CameraX = Mouse.X/Tile_Size_W + CameraX
			CameraY = Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}
	}
	if Keyboard.GetKey("3") == 1 {
		if releasedZoom == true {
			Tile_Size_H = 8
			Tile_Size_W = 8
			CameraX = Mouse.X/Tile_Size_W + CameraX
			CameraY = Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}
	}

	if Keyboard.GetKey("4") == 1 {
		if releasedZoom == true {
			Tile_Size_H = 4
			Tile_Size_W = 4
			CameraX = 0
			CameraY = 0
			releasedZoom = false
		}
	}

	if Keyboard.GetKey("4") == 0 && Keyboard.GetKey("1") == 0 && Keyboard.GetKey("2") == 0 && Keyboard.GetKey("3") == 0 && Keyboard.GetKey("4") == 0 {
		releasedZoom = true
	}

	viewWidth := World_W / Tile_Size_W
	viewHeight := World_H / Tile_Size_H

	view := level.GetView(CameraX, CameraY, viewWidth, viewHeight, false, false)

	renderer.Clear()
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := int32(x * Tile_Size_W)
			tY := int32(y * Tile_Size_H)
			tile := view[x][y]

			drawSprite(tX, tY, tile.SpriteX, tile.SpriteY, 255, 255, 255, worldTexture) //Tile itself

			if tile == nil {
				drawSprite(tX, tY, 0, 112, 255, 255, 255, worldTexture) //Empty space
			} else {
				//Draw entity on tile.
				entity := level.GetEntityAt(tile.X, tile.Y)
				if entity != nil {
					if entity.HasComponent("AppearanceComponent") {
						ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
						dir := 0
						if entity.HasComponent("DirectionComponent") {
							dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
							dir = dc.Direction
						}
						drawSprite(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY, ac.R, ac.G, ac.B, characterTexture) //Entity

						//Temp select code
						if pX == tile.X && pY == tile.Y && Mouse.Clicked {
							for _, entity := range level.Entities {
								if entity.HasComponent("SelectedComponent") {
									entity.RemoveComponent("SelectedComponent")
								}
							}
							entity.AddComponent(component.SelectedComponent{})
						}
						if entity.HasComponent("SelectedComponent") {
							drawSprite(tX, tY, 112, 128, 255, 255, 255, uiTexture)
							if entity.HasComponent("DescriptionComponent") {
								dc := entity.GetComponent("DescriptionComponent").(*component.DescriptionComponent)
								drawText(World_W, 10, dc.Name)
							}
							if entity.HasComponent("GoblinAIComponent") {
								gc := entity.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)

								drawText(World_W, 50, gc.State)
								drawText(World_W, 75, "Energy:"+strconv.Itoa(gc.Energy))
							}

						}
					}
				}
				if pX == tile.X && pY == tile.Y {
					drawSprite(tX, tY, 128, 128, 255, 255, 255, uiTexture) //Cursor?
				}
			}
		}
	}

	renderer.Present()

	return level
}

func drawSprite(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	renderer.Copy(texture, &src, &dst)
}

func drawText(x int32, y int32, text string) {

	var solidTexture *sdl.Texture
	var err error

	var solidSurface *sdl.Surface
	if solidSurface, err = font.RenderUTF8BlendedWrapped(text, sdl.Color{255, 255, 255, 255}, 200); err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
		return
	}
	if solidTexture, err = renderer.CreateTextureFromSurface(solidSurface); err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
		return
	}
	solidSurface.Free()
	_, _, w, h, err := solidTexture.Query()
	dst := sdl.Rect{X: x, Y: y, W: w, H: h}
	renderer.Copy(solidTexture, nil, &dst)
	solidTexture.Destroy()
}
