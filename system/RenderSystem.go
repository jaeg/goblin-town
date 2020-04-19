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
const Window_W = 992
const Window_H = 576
const World_W = 800
const World_H = 576
const MiniMap_X = World_W
const MiniMap_Y = 200

type RenderSystem struct {
}

var renderer *sdl.Renderer
var characterTexture *sdl.Texture
var worldTexture *sdl.Texture
var uiTexture *sdl.Texture
var window *sdl.Window
var font *ttf.Font
var miniMapTexture *sdl.Texture

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

	//Generate mini map

	sdl.ShowCursor(0)
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
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: World_W, H: World_H})

	//Draw menu
	for x := World_W; x < Window_W; x += 16 {
		for y := 0; y < Window_H; y += 16 {
			var sX int32
			var sY int32
			sX = 128
			sY = 16
			//Left Top
			if x == World_W && y == 0 {
				sY = 0
				sX = 144
			} else if x == Window_W-16 && y == 0 { //Right top
				sY = 0
				sX = 176
			} else if x == World_W && y == Window_H-16 { //Left bottom
				sY = 0
				sX = 144
			} else if x == Window_W-16 && y == Window_H-16 { //Right bottom
				sY = 0
				sX = 144
			}
			drawSprite(int32(x), int32(y), sX, sY, 255, 255, 255, uiTexture)
		}
	}

	if miniMapTexture == nil {
		CreateMiniMap(level)
		fmt.Println("Create mini map")
	} else {
		_, _, w, h, _ := miniMapTexture.Query()
		src := sdl.Rect{X: 0, Y: 0, W: w, H: h}
		dst := sdl.Rect{X: MiniMap_X, Y: MiniMap_Y, W: 200, H: 200}
		renderer.Copy(miniMapTexture, &src, &dst)

		renderer.DrawRect(&sdl.Rect{X: int32(MiniMap_X + CameraX), Y: int32(MiniMap_Y + CameraY), W: int32(len(view)), H: int32(len(view[0]))})
	}

	//Draw world
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := int32(x * Tile_Size_W)
			tY := int32(y * Tile_Size_H)
			tile := view[x][y]

			if tile == nil {
				drawSprite(tX, tY, 0, 112, 255, 255, 255, worldTexture) //Empty space
			} else {
				//For tiles we want the higher tiles to appear whiter to convey depth.  The background is white
				//so changing the alpha does this for us.
				alpha := 255
				if tile.Elevation > 0 {
					alpha = 255 - tile.Elevation*4
					// It doesn't make sense to start off by fading out the bottom of the mountain.
					if tile.Elevation == 2 {
						alpha = 255
					}
					if alpha > 255 {
						alpha = 255
					}
				}

				depth := 255
				if tile.Elevation < 0 {
					depth = 255 + tile.Elevation*20
				}

				drawSpriteEx(tX, tY, tile.SpriteX, tile.SpriteY, uint8(depth), uint8(depth), uint8(depth), uint8(alpha), worldTexture) //Tile itself

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
						if entity.HasComponent("DeadComponent") {
							drawSpriteUpsideDown(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY, ac.R, ac.G, ac.B, characterTexture) //Entity
						} else {
							drawSprite(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY, ac.R, ac.G, ac.B, characterTexture) //Entity
						}

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

							if entity.HasComponent("HealthComponent") {
								gc := entity.GetComponent("HealthComponent").(*component.HealthComponent)

								drawText(World_W, 85, "Hp:"+strconv.Itoa(gc.Health))
							}

							if entity.HasComponent("DeadComponent") {
								drawText(World_W, 0, "Dead")
							}

						}
					}
				}
				var cursorY int32
				cursorY = 128
				if Mouse.Clicked {
					cursorY = 144
					if Mouse.X > MiniMap_X && Mouse.X < MiniMap_X+200 && Mouse.Y > MiniMap_Y && Mouse.Y < MiniMap_Y+200 {
						newCameraX := Mouse.X - MiniMap_X
						newCameraY := Mouse.Y - MiniMap_Y
						if newCameraX+viewWidth <= level.Width {
							CameraX = newCameraX
						} else {
							CameraX = level.Width - viewWidth
						}

						if newCameraY+viewHeight <= level.Height {
							CameraY = newCameraY
						} else {
							CameraY = level.Height - viewHeight
						}

					}
				}
				if pX == tile.X && pY == tile.Y {
					drawSprite(tX, tY, 128, cursorY, 255, 255, 255, uiTexture) //Cursor?
				}

				if Mouse.X > World_W {
					drawSprite(int32(Mouse.X), int32(Mouse.Y), 64, cursorY, 255, 255, 255, uiTexture) //Cursor?
				}
			}
		}
	}

	renderer.Present()

	return level
}

func drawSprite(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	texture.SetAlphaMod(255)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	renderer.Copy(texture, &src, &dst)

}

func drawSpriteEx(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, a uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	texture.SetAlphaMod(a)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	renderer.Copy(texture, &src, &dst)

}

func drawSpriteUpsideDown(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	renderer.CopyEx(texture, &src, &dst, 0, nil, sdl.FLIP_VERTICAL)

}

func drawText(x int32, y int32, text string) {

	var solidTexture *sdl.Texture
	var err error

	var solidSurface *sdl.Surface
	if solidSurface, err = font.RenderUTF8BlendedWrapped(text, sdl.Color{255, 255, 255, 255}, 192); err != nil {
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

func CreateMiniMap(level *world.Level) {
	image, err := img.Load("tiny_dungeon_world.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s\n", err)
		return
	}

	surface, err := sdl.CreateRGBSurface(0, int32(level.Width), int32(level.Height), 16, 0, 0, 0, 0)

	if err != nil {
		fmt.Printf("Failed to create minimap surface: %s\n", err)
		return
	}
	//Draw minimap
	for x := 0; x < level.Width; x++ {
		for y := 0; y < level.Height; y++ {
			tX := int32(x)
			tY := int32(y)
			tile := level.GetTileAt(x, y)

			src := &sdl.Rect{X: tile.SpriteX, Y: tile.SpriteY, W: Sprite_Size_W, H: Sprite_Size_H}
			dst := &sdl.Rect{X: tX, Y: tY, W: 1, H: 1}
			err = image.Blit(src, surface, dst)
			if err != nil {
				fmt.Printf("Failed to create minimap surface: %s\n", err)
				return
			}
		}
	}

	if miniMapTexture, err = renderer.CreateTextureFromSurface(surface); err != nil {
		fmt.Printf("Failed to create minimap texture: %s\n", err)
		return
	}
	image.Free()
	surface.Free()
}
