package system

import (
	"github.com/jaeg/goblin-town/component"
	"github.com/jaeg/goblin-town/entity"
	entityFactory "github.com/jaeg/goblin-town/entity"
	"github.com/jaeg/goblin-town/world"

	"github.com/beefsack/go-astar"
)

var GoblinTorch_X = 0
var GoblinTorch_Y = 0

// GoblinAISystem Manages the goblin's ai in the simulation
type GoblinAISystem struct {
}

// Update Main update function of the system.
func (s GoblinAISystem) Update(level *world.Level, entity *entity.Entity) *world.Level {
	if entity.HasComponent("GoblinAIComponent") {
		if entity.HasComponent("MyTurnComponent") {
			aic := entity.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)

			deltaX := 0
			deltaY := 0
			if entity.HasComponent("HealthComponent") {
				hc := entity.GetComponent("HealthComponent").(*component.HealthComponent)
				if hc.Health <= 0 {
					if entity.HasComponent("FoodComponent") {
						entity.RemoveComponent("GoblinAIComponent")
					} else {
						entity.AddComponent(&component.DeadComponent{})
					}
					return level
				}
			}

			switch aic.State {
			case "wander":
				deltaX = getRandom(-1, 2)
				deltaY = 0
				if deltaX == 0 {
					deltaY = getRandom(-1, 2)
				}

				if aic.Energy < aic.HungerThreshold {
					aic.State = "search"
				}

				//count goblins near me
				goblinsNearby := 0
				nearby := level.GetEntitiesAround(pc.GetX(), pc.GetY(), aic.SightRange, aic.SightRange)
				for e := range nearby {
					if nearby[e].HasComponent("GoblinAIComponent") && !nearby[e].HasComponent("DeadComponent") {
						goblinsNearby++
					}
				}

				if goblinsNearby < aic.SocialThreshold {
					aic.State = "findfriends"
				} else {
					//count goblins near me for mating purposes
					goblinsNearby := 0
					emptyX := -1
					emptyY := -1
					//Look one square around and count them goblins
					for x := pc.GetX() - 1; x < pc.GetX()+1; x++ {
						for y := pc.GetY() - 1; y < pc.GetY()+1; y++ {
							tile := level.GetTileAt(x, y)
							if tile != nil {
								entityHit := level.GetSolidEntityAt(x, y)
								if entityHit != nil {
									if entityHit.HasComponent("GoblinAIComponent") && !entityHit.HasComponent("DeadComponent") {
										goblinsNearby++
									}
								} else {
									//There's no solid entity in this square.  If it's not a blocked tile it's possible to birth a goblin here.
									if tile.Type != 2 && tile.Type != 4 {
										emptyX = x
										emptyY = y
									}
								}
							}
						}
					}
					//Make a new goblin
					if emptyX != -1 && emptyY != -1 && goblinsNearby >= aic.MateThreshold && aic.Energy > aic.HungerThreshold {
						goblin, err := entityFactory.Create("goblin", emptyX, emptyY)
						newAic := goblin.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)
						energy := aic.Energy / 2
						aic.Energy = energy
						newAic.Energy = energy
						if err == nil {
							level.AddEntity(goblin)
						}
					}
				}
			case "findfriends":
				nearby := level.GetEntitiesAround(pc.GetX(), pc.GetY(), aic.SightRange, aic.SightRange)
				// If goblins don't have anyone nearby to go to, head to the torch for saftey.
				if len(nearby) < 2 {
					aic.TargetX = GoblinTorch_X
					aic.TargetY = GoblinTorch_Y
					aic.State = "approach"
				} else {
					// Find a nearby goblin and head towards it.
					for e := range nearby {
						if nearby[e].HasComponent("GoblinAIComponent") && !nearby[e].HasComponent("DeadComponent") {
							goblinPC := nearby[e].GetComponent("PositionComponent").(*component.PositionComponent)
							aic.TargetX = goblinPC.GetX()
							aic.TargetY = goblinPC.GetY()
							aic.State = "approach"
							break
						}
					}
				}
			case "approach":
				steps, _, _ := astar.Path(level.GetTileAt(pc.GetX(), pc.GetY()), level.GetTileAt(aic.TargetX, aic.TargetY))
				if len(steps) > 0 {
					t := steps[0].(*world.Tile)
					if pc.GetX() < t.X {
						deltaX = 1
					}

					if pc.GetX() > t.X {
						deltaX = -1
					}

					if pc.GetY() < t.Y {
						deltaY = 1
					}

					if pc.GetY() > t.Y {
						deltaY = -1
					}

				}

				if deltaX == 0 && deltaY == 0 {
					aic.State = "wander"
				}
			case "search":
				//Scan around for food to the best my vision allows me.
				nearby := level.GetEntitiesAround(pc.GetX(), pc.GetY(), aic.SightRange, aic.SightRange)

				for e := range nearby {
					if nearby[e].HasComponent("FoodComponent") && !nearby[e].HasComponent("DeadComponent") {
						foodPC := nearby[e].GetComponent("PositionComponent").(*component.PositionComponent)
						aic.TargetX = foodPC.GetX()
						aic.TargetY = foodPC.GetY()
						aic.State = "approach"
						break
					}
				}

				//Found nothing, wander
				if deltaX == 0 && deltaY == 0 {
					deltaX = getRandom(-1, 2)
					deltaY = 0
					if deltaX == 0 {
						deltaY = getRandom(-1, 2)
					}
				}
			}

			//Move
			if move(entity, level, deltaX, deltaY) {
				entityHit := level.GetSolidEntityAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
				if entityHit != nil {
					if entityHit != entity && entityHit.Blueprint != "goblin" {
						hit(entity, entityHit)
						if eat(entity, entityHit) {
							aic.Energy += 2
						}
					}
				}
			}
			face(entity, deltaX, deltaY)

			//Every thing we just did costed energy.
			aic.Energy--
			if aic.Energy <= 0 {
				entity.AddComponent(&component.DeadComponent{})
			}
		}
	}

	return level
}
