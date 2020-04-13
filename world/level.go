package world

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jaeg/goblin-town/entity"

	"github.com/jaeg/goblin-town/component"

	"github.com/aquilax/go-perlin"
)

const (
	alpha = 6.
	beta  = 5.
	n     = 5
)

// Level .
type Level struct {
	data                  [][]Tile
	Entities              []*entity.Entity
	width, height         int
	id                    int
	left, right, up, down int
	theme                 string
}

//Tile Types
//1 - open
//2 - solid
//3 - stairs [level id, to x, to y]
//4 - water
type Tile struct {
	Type      int
	SpriteX   int32
	SpriteY   int32
	Data      []int
	X         int
	Y         int
	Elevation int
}

func newLevel(width int, height int) (level *Level) {
	level = &Level{width: width, height: height, left: -1, right: -1, up: -1, down: -1}

	data := make([][]Tile, width, height)
	for x := 0; x < width; x++ {
		col := []Tile{}
		for y := 0; y < height; y++ {
			col = append(col, Tile{Type: 1, X: x, Y: y, SpriteX: 16, SpriteY: 128})
		}
		data[x] = append(data[x], col...)
	}

	level.data = data
	return
}

func NewOverworldSection(width int, height int) (level *Level) {
	fmt.Println("Creating new random level")
	level = newLevel(width, height)

	p := perlin.NewPerlin(alpha, beta, n, time.Now().UnixNano())
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			tile := level.GetTileAt(x, y)
			value := int(p.Noise2D(float64(x)/100, float64(y)/100) * 10)
			tile.Elevation = value
			if value == -1 {
				tile.Type = 0
				tile.SpriteX = 176
				tile.SpriteY = 80
			}

			//grass
			if value > -1 {
				tile.Type = 0
				tile.SpriteX = 128
				tile.SpriteY = 80
			}

			if value >= 2 {
				tile.Type = 0
				tile.SpriteX = 0
				tile.SpriteY = 80
			}
		}
	}
	/*
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if rand.Intn(5) == 0 {
					level.GetTileAt(x, y).SpriteX = 144
				} else if rand.Intn(50) == 0 {
					level.GetTileAt(x, y).SpriteX = 192
				}
			}
		}

		//Generate flower formations
		for i := 0; i < 50; i++ {
			x := getRandom(1, width)
			y := getRandom(1, height)

			level.createCluster(x, y, 10, 176, 80, false, false)
		}

		//Generate Water
		for i := 0; i < 1000; i++ {
			x := getRandom(1, width)
			y := getRandom(1, height)

			level.createCluster(x, y, getRandom(20, 300), 16, 128, false, true)
		}

		//Generate Mountains
		for i := 0; i < 3; i++ {
			x := getRandom(1, width)
			y := getRandom(1, height)

			level.createCluster(x, y, 2000, 0, 80, true, false)
		}
	*/
	return
}

func (level *Level) GetTileAt(x int, y int) (tile *Tile) {
	if x < level.width && y < level.height && x >= 0 && y >= 0 {
		tile = &level.data[x][y]
	}
	return
}

//Get's the view frustum with the player in the center
func (level *Level) GetView(aX int, aY int, width int, height int, blind bool, centered bool) (data [][]*Tile) {
	left := aX - width/2
	right := aX + width/2
	up := aY - height/2
	down := aY + height/2

	if centered == false {
		left = aX
		right = aX + width - 1
		up = aY
		down = aY + height
	}

	data = make([][]*Tile, width+1-width%2)

	cX := 0
	for x := left; x <= right; x++ {
		col := []*Tile{}
		for y := up; y <= down; y++ {
			currentTile := level.GetTileAt(x, y)
			if blind {
				if y < aY-height/4 || y > aY+height/4 || x > aX+width/4 || x < aX-width/4 {
					currentTile = nil
				}
			}
			col = append(col, currentTile)
		}
		data[cX] = append(data[cX], col...)
		cX++
	}
	return
}

func (level *Level) GetEntitiesAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	entitiesLen := len(level.Entities)

	for i := 0; i < entitiesLen; i++ {
		entity := level.Entities[i]

		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X >= left && pc.X <= right && pc.Y >= up && pc.Y <= down {
				entities = append(entities, entity)
			}
		}
	}

	return
}

func (level *Level) GetPlayersAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	entitiesLen := len(level.Entities)

	for i := 0; i < entitiesLen; i++ {
		entity := level.Entities[i]

		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X >= left && pc.X <= right && pc.Y >= up && pc.Y <= down {
				if entity.HasComponent("PlayerComponent") {
					entities = append(entities, entity)
				}
			}
		}
	}

	return
}

func (level *Level) GetEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X == x && pc.Y == y {
				return
			}
		}
	}
	entity = nil
	return
}

func (level *Level) GetSolidEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			if entity.HasComponent("SolidComponent") {
				pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
				if pc.X == x && pc.Y == y {
					return
				}
			}
		}
	}
	entity = nil
	return
}

func (level *Level) GetInteractableEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			if entity.HasComponent("InteractComponent") {
				pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
				if pc.X == x && pc.Y == y {
					return
				}
			}
		}
	}
	entity = nil
	return
}

func (level *Level) AddEntity(entity *entity.Entity) {
	level.Entities = append(level.Entities, entity)
}

func (level *Level) RemoveEntity(entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		if level.Entities[i] == entity {
			level.Entities = append(level.Entities[:i], level.Entities[i+1:]...)

		}
	}
}

func (level *Level) createCluster(x int, y int, size int, spriteX int32, spriteY int32, solid bool, water bool) {
	for i := 0; i < size; i++ {
		n := getRandom(1, 6)
		e := getRandom(1, 6)
		s := getRandom(1, 6)
		w := getRandom(1, 6)

		if n == 1 {
			x += 1
		}

		if s == 1 {
			x--
		}

		if e == 1 {
			y++
		}

		if w == 1 {
			y--
		}

		if level.GetTileAt(x, y) != nil {
			tile := level.GetTileAt(x, y)
			tile.SpriteX = spriteX
			tile.SpriteY = spriteY
			if solid {
				tile.Type = 2
			} else if water {
				tile.Type = 4
			} else {
				tile.Type = 1
			}

		}
	}
}

func getRandom(low int, high int) int {
	if low == high {
		return low
	}
	return (rand.Intn((high - low))) + low
}

func Sgn(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

//Ported from http://www.roguebasin.com/index.php?title=Simple_Line_of_Sight
func los(pX int, pY int, tX int, tY int, level *Level) bool {
	deltaX := pX - tX
	deltaY := pY - tY

	absDeltaX := math.Abs(float64(deltaX))
	absDeltaY := math.Abs(float64(deltaY))

	signX := Sgn(deltaX)
	signY := Sgn(deltaY)

	if absDeltaX > absDeltaY {
		t := absDeltaY*2 - absDeltaX
		for {
			if t >= 0 {
				tY += signY
				t -= absDeltaX * 2
			}

			tX += signX
			t += absDeltaY * 2

			if tX == pX && tY == pY {
				return true
			}
			if level.GetTileAt(tX, tY).Type == 2 {
				break
			}
		}
		return false
	}

	t := absDeltaX*2 - absDeltaY

	for {
		if t >= 0 {
			tX += signX
			t -= absDeltaY * 2
		}
		tY += signY
		t += absDeltaX * 2
		if tX == pX && tY == pY {
			return true
		}

		if level.GetTileAt(tX, tY).Type == 2 {
			break
		}
	}

	return false

}

func distance(x1 int, y1 int, x2 int, y2 int) int {
	var dy int
	if y1 > y2 {
		dy = y1 - y2
	} else {
		dy = (y2 - y1)
	}

	var dx int
	if x1 > x2 {
		dx = x1 - x2
	} else {
		dx = x2 - x1
	}

	var d int
	if dy > dx {
		d = dy + (dx >> 1)
	} else {
		d = dx + (dy >> 1)
	}

	return d
}
