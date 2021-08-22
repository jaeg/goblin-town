package render

import (
	"log"
	"strconv"

	"github.com/jaeg/goblin-town/component"
	"github.com/jaeg/goblin-town/system"
	"github.com/jaeg/goblin-town/world"

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

const MiniMap_X = World_W + 6
const MiniMap_Y = 200
const MiniMap_W = 180
const MiniMap_H = 180

type Renderer struct {
	renderer         *sdl.Renderer
	characterTexture *sdl.Texture
	fxTexture        *sdl.Texture
	worldTexture     *sdl.Texture
	uiTexture        *sdl.Texture
	window           *sdl.Window
	font             *ttf.Font
	miniMapTexture   *sdl.Texture
}

var CameraX = 0
var CameraY = 0

var Zoom = 1

var releasedZoom = true

var Beat = 0

func (s *Renderer) Init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	s.window, err = sdl.CreateWindow("Tiles", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		Window_W, Window_H, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Printf("Failed to create window: %s\n", err)
		return
	}

	if err = ttf.Init(); err != nil {
		log.Printf("Failed to initialize TTF: %s\n", err)
		return
	}

	s.renderer, err = sdl.CreateRenderer(s.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Printf("Failed to create renderer: %s\n", err)
		return
	}

	image, err := img.Load("tiny_dungeon_monsters.png")
	if err != nil {
		log.Printf("Failed to load BMP: %s\n", err)
		return
	}

	s.characterTexture, err = s.renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Printf("Failed to create texture: %s\n", err)
		return
	}

	image.Free()

	image, err = img.Load("tiny_dungeon_fx.png")
	if err != nil {
		log.Printf("Failed to load BMP: %s\n", err)
		return
	}

	s.fxTexture, err = s.renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Printf("Failed to create texture: %s\n", err)
		return
	}

	image.Free()

	image, err = img.Load("tiny_dungeon_world.png")
	if err != nil {
		log.Printf("Failed to load BMP: %s\n", err)
		return
	}

	s.worldTexture, err = s.renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Printf("Failed to create texture: %s\n", err)
		return
	}

	image.Free()

	image, err = img.Load("tiny_dungeon_interface.png")
	if err != nil {
		log.Printf("Failed to load BMP: %s\n", err)
		return
	}

	s.uiTexture, err = s.renderer.CreateTextureFromSurface(image)
	if err != nil {
		log.Printf("Failed to create texture: %s\n", err)
		return
	}

	if s.font, err = ttf.OpenFont("Roboto-Regular.ttf", 30); err != nil {
		log.Printf("Failed to open font: %s\n", err)
		return
	}

	//Generate mini map

	sdl.ShowCursor(0)
}

func (s *Renderer) Cleanup() {
	sdl.Quit()
	s.window.Destroy()
	s.renderer.Destroy()
	s.worldTexture.Destroy()
	s.characterTexture.Destroy()
	s.uiTexture.Destroy()

}

// Renderer Update
func (s *Renderer) Update(level *world.Level) *world.Level {
	pX := system.Mouse.X/Tile_Size_W + CameraX
	pY := system.Mouse.Y/Tile_Size_H + CameraY

	if system.Keyboard.GetKey("a") == 1 && CameraX > 0 {
		CameraX--
	}
	if system.Keyboard.GetKey("d") == 1 && CameraX < level.Width-World_W/Tile_Size_W-1 {
		CameraX++
	}
	if system.Keyboard.GetKey("w") == 1 && CameraY > 0 {
		CameraY--
	}
	if system.Keyboard.GetKey("s") == 1 && CameraY < level.Height-World_H/Tile_Size_H-1 {
		CameraY++
	}

	if system.Keyboard.GetKey("1") == 1 {
		if releasedZoom {
			Tile_Size_H = 32
			Tile_Size_W = 32
			CameraX = system.Mouse.X/Tile_Size_W + CameraX
			CameraY = system.Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}

	}
	if system.Keyboard.GetKey("2") == 1 {
		if releasedZoom {
			Tile_Size_H = 16
			Tile_Size_W = 16
			CameraX = system.Mouse.X/Tile_Size_W + CameraX
			CameraY = system.Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}
	}
	if system.Keyboard.GetKey("3") == 1 {
		if releasedZoom {
			Tile_Size_H = 8
			Tile_Size_W = 8
			CameraX = system.Mouse.X/Tile_Size_W + CameraX
			CameraY = system.Mouse.Y/Tile_Size_H + CameraY

			if CameraY > level.Height-World_H/Tile_Size_H-1 {
				CameraY = level.Height - World_H/Tile_Size_H - 1
			}

			if CameraX > level.Width-World_W/Tile_Size_W-1 {
				CameraX = level.Width - World_W/Tile_Size_W - 1
			}
			releasedZoom = false
		}
	}

	if system.Keyboard.GetKey("4") == 1 {
		if releasedZoom {
			Tile_Size_H = 4
			Tile_Size_W = 4
			CameraX = 0
			CameraY = 0
			releasedZoom = false
		}
	}

	if system.Keyboard.GetKey("4") == 0 && system.Keyboard.GetKey("1") == 0 && system.Keyboard.GetKey("2") == 0 && system.Keyboard.GetKey("3") == 0 && system.Keyboard.GetKey("4") == 0 {
		releasedZoom = true
	}

	viewWidth := World_W / Tile_Size_W
	viewHeight := World_H / Tile_Size_H

	view := level.GetView(CameraX, CameraY, viewWidth, viewHeight, false, false)

	s.renderer.Clear()
	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: World_W, H: World_H})

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
			s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
		}
	}

	if s.miniMapTexture == nil {
		s.CreateMiniMap(level)
		log.Println("Create mini map")
	} else {
		_, _, w, h, _ := s.miniMapTexture.Query()
		src := sdl.Rect{X: 0, Y: 0, W: w, H: h}
		dst := sdl.Rect{X: MiniMap_X, Y: MiniMap_Y, W: MiniMap_W, H: MiniMap_H}
		s.renderer.Copy(s.miniMapTexture, &src, &dst)

		scale := float64(MiniMap_W) / float64(level.Width)

		scaledX := MiniMap_X + float64(CameraX)*scale
		scaledY := MiniMap_Y + float64(CameraY)*scale
		//Box
		s.renderer.DrawRect(&sdl.Rect{X: int32(scaledX), Y: int32(scaledY), W: int32(float64(len(view)) * scale), H: int32(float64(len(view[0])) * scale)})

		//Torch
		scaledX = MiniMap_X + float64(system.GoblinTorch_X)*scale
		scaledY = MiniMap_Y + float64(system.GoblinTorch_Y)*scale
		s.renderer.SetDrawColor(255, 0, 0, 255)
		s.renderer.DrawRect(&sdl.Rect{X: int32(scaledX - 2), Y: int32(scaledY - 2), W: 4, H: 4})
		s.renderer.SetDrawColor(255, 255, 255, 255)

	}

	//Draw world
	for x := 0; x < len(view); x++ {
		for y := 0; y < len(view[x]); y++ {
			tX := int32(x * Tile_Size_W)
			tY := int32(y * Tile_Size_H)
			tile := view[x][y]

			if tile == nil {
				s.drawSprite(tX, tY, 0, 112, 255, 255, 255, s.worldTexture) //Empty space
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

				s.drawSpriteEx(tX, tY, tile.SpriteX, tile.SpriteY, int32(Tile_Size_W), int32(Tile_Size_H), uint8(depth), uint8(depth), uint8(depth), uint8(alpha), s.worldTexture) //Tile itself

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

						if entity.HasComponent("InanimateComponent") {
							s.drawSprite(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY, ac.R, ac.G, ac.B, s.worldTexture)
						} else {
							if entity.HasComponent("DeadComponent") {
								s.drawSpriteUpsideDown(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY, ac.R, ac.G, ac.B, s.characterTexture) //Entity
							} else {
								//Entity
								s.drawSprite(tX, tY, ac.SpriteX+(int32(dir)*Sprite_Size_W), ac.SpriteY+(int32(Beat)*Sprite_Size_H), ac.R, ac.G, ac.B, s.characterTexture)
								//Draw FX
								if entity.HasComponent("AttackComponent") {
									attackC := entity.GetComponent("AttackComponent").(*component.AttackComponent)
									if attackC.Frame == 3 {
										entity.RemoveComponent("AttackComponent")
									} else {
										s.drawSprite(tX, tY, int32(attackC.SpriteX)+(int32(attackC.Frame)*Sprite_Size_W), int32(attackC.SpriteY), 255, 255, 255, s.fxTexture)
										attackC.Frame++
									}
								}
							}
						}

						//Temp select code
						if pX == tile.X && pY == tile.Y && system.Mouse.Clicked {
							for _, entity := range level.Entities {
								if entity.HasComponent("SelectedComponent") {
									entity.RemoveComponent("SelectedComponent")
								}
							}
							entity.AddComponent(component.SelectedComponent{})
						}
						if entity.HasComponent("SelectedComponent") {
							s.drawSprite(tX, tY, 112, 128, 255, 255, 255, s.uiTexture)

							if entity.HasComponent("DescriptionComponent") {
								dc := entity.GetComponent("DescriptionComponent").(*component.DescriptionComponent)
								s.drawText(World_W, 10, dc.Name)
							}
							if entity.HasComponent("GoblinAIComponent") {
								gc := entity.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)

								s.drawText(World_W, 50, gc.State)
								s.drawText(World_W, 75, "Energy:"+strconv.Itoa(gc.Energy))
							}

							if entity.HasComponent("HealthComponent") {
								gc := entity.GetComponent("HealthComponent").(*component.HealthComponent)

								s.drawText(World_W, 85, "Hp:"+strconv.Itoa(gc.Health))
							}

							if entity.HasComponent("DeadComponent") {
								s.drawText(World_W, 0, "Dead")
							}

						}
					}
				}

				//Torch
				if system.GoblinTorch_X == tile.X && system.GoblinTorch_Y == tile.Y {
					s.drawSprite(tX, tY, 80+int32(Sprite_Size_W*Beat), 192, 255, 255, 255, s.worldTexture)
				}

				//Cursor stuff
				var cursorY int32
				cursorY = 128
				if system.Mouse.Clicked {
					cursorY = 144
					if system.Mouse.X > MiniMap_X && system.Mouse.X < MiniMap_X+MiniMap_W && system.Mouse.Y > MiniMap_Y && system.Mouse.Y < MiniMap_Y+MiniMap_H {
						scale := float64(MiniMap_W) / float64(level.Width)
						scaledX := int(float64(system.Mouse.X-MiniMap_X) / scale)
						scaledY := int(float64(system.Mouse.Y-MiniMap_Y) / scale)
						PlaceCamera(int(scaledX), int(scaledY), level)

						PlaceCamera(scaledX, scaledY, level)
					}
				}
				if pX == tile.X && pY == tile.Y {
					s.drawSprite(tX, tY, 128, cursorY, 255, 255, 255, s.uiTexture) //Cursor?
				}

				if system.Mouse.X > World_W {
					s.drawSprite(int32(system.Mouse.X), int32(system.Mouse.Y), 64, cursorY, 255, 255, 255, s.uiTexture) //Cursor?
				}
			}
		}
	}

	//Render the day/night

	/*
		Get's brighter between 5am and 8am.
		Get's darker between 7pm and 9pm.
		Stays dark between 9pm and 5am.
	*/
	//Dawn - Get's brighter between 5am and 8am.
	if level.Hour >= 5 && level.Hour < 8 {
		alpha := 125 - 10*level.Hour
		s.renderer.SetDrawColor(0, 0, 0, uint8(alpha))
		s.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		s.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: World_W, H: World_H})
		s.renderer.SetDrawColor(255, 255, 255, 255)
	} else if level.Hour > 17 && level.Hour <= 21 { //Dusk - Get's darker between 7pm and 9pm.
		alpha := 0 + 25*(level.Hour-16)
		s.renderer.SetDrawColor(0, 0, 0, uint8(alpha))
		s.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		s.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: World_W, H: World_H})
		s.renderer.SetDrawColor(255, 255, 255, 255)
	} else if level.Hour > 21 || (level.Hour >= 0 && level.Hour < 5) {
		s.renderer.SetDrawColor(0, 0, 0, 125)
		s.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		s.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: World_W, H: World_H})
		s.renderer.SetDrawColor(255, 255, 255, 255)
	}

	s.renderer.Present()

	return level
}

func (s *Renderer) drawSprite(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	texture.SetAlphaMod(255)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	s.renderer.Copy(texture, &src, &dst)
}

func (s *Renderer) drawSpriteEx(x int32, y int32, sx int32, sy int32, w int32, h int32, r uint8, g uint8, b uint8, a uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	texture.SetAlphaMod(a)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: w, H: h}
	s.renderer.Copy(texture, &src, &dst)

}

func (s *Renderer) drawSpriteUpsideDown(x int32, y int32, sx int32, sy int32, r uint8, g uint8, b uint8, texture *sdl.Texture) {
	texture.SetColorMod(r, g, b)
	src := sdl.Rect{X: sx, Y: sy, W: Sprite_Size_W, H: Sprite_Size_H}
	dst := sdl.Rect{X: x, Y: y, W: int32(Tile_Size_W), H: int32(Tile_Size_H)}
	s.renderer.CopyEx(texture, &src, &dst, 0, nil, sdl.FLIP_VERTICAL)

}

func (s *Renderer) drawText(x int32, y int32, text string) {

	var solidTexture *sdl.Texture
	var err error

	var solidSurface *sdl.Surface
	if solidSurface, err = s.font.RenderUTF8BlendedWrapped(text, sdl.Color{R: 255, G: 255, B: 255, A: 255}, 192); err != nil {
		log.Printf("Failed to render text: %s\n", err)
		return
	}

	if solidTexture, err = s.renderer.CreateTextureFromSurface(solidSurface); err != nil {
		log.Printf("Failed to create texture: %s\n", err)
		return
	}
	solidSurface.Free()
	_, _, w, h, err := solidTexture.Query()
	if err != nil {
		log.Printf("Error querying texture")
	}
	dst := sdl.Rect{X: x, Y: y, W: w, H: h}
	s.renderer.Copy(solidTexture, nil, &dst)
	solidTexture.Destroy()
}

func (s *Renderer) CreateMiniMap(level *world.Level) {
	image, err := img.Load("tiny_dungeon_world.png")
	if err != nil {
		log.Printf("Failed to load BMP: %s\n", err)
		return
	}

	surface, err := sdl.CreateRGBSurface(0, int32(level.Width), int32(level.Height), 16, 0, 0, 0, 0)

	if err != nil {
		log.Printf("Failed to create minimap surface: %s\n", err)
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
				log.Printf("Failed to create minimap surface: %s\n", err)
				return
			}
		}
	}

	if s.miniMapTexture, err = s.renderer.CreateTextureFromSurface(surface); err != nil {
		log.Printf("Failed to create minimap texture: %s\n", err)
		return
	}
	image.Free()
	surface.Free()
}

func PlaceCamera(x int, y int, level *world.Level) {
	newCameraX := x
	newCameraY := y
	viewWidth := World_W / Tile_Size_W
	viewHeight := World_H / Tile_Size_H
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

func CenterCamera(x int, y int, level *world.Level) {

	viewWidth := World_W / Tile_Size_W
	viewHeight := World_H / Tile_Size_H
	newCameraX := x - viewWidth/2
	newCameraY := y - viewHeight/2
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
