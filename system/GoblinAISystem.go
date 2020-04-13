package system

import (
	"github.com/jaeg/goblin-town/world"

	entityFactory "github.com/jaeg/goblin-town/entity"

	"github.com/jaeg/goblin-town/entity"

	"github.com/jaeg/goblin-town/component"
)

type GoblinAISystem struct {
}

// GoblinAISystem .
func (s GoblinAISystem) Update(level *world.Level, entity *entity.Entity) *world.Level {
	if entity.HasComponent("GoblinAIComponent") {
		if entity.HasComponent("MyTurnComponent") {
			aic := entity.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)

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

				for e := range nearby {
					if nearby[e].HasComponent("GoblinAIComponent") && !nearby[e].HasComponent("DeadComponent") {
						goblinPC := nearby[e].GetComponent("PositionComponent").(*component.PositionComponent)
						aic.TargetX = goblinPC.GetX()
						aic.TargetY = goblinPC.GetY()
						aic.State = "approach"
						break
					}
				}
			case "approach":
				if pc.GetX() < aic.TargetX {
					deltaX = 1
				}

				if pc.GetX() > aic.TargetX {
					deltaX = -1
				}

				if pc.GetY() < aic.TargetY {
					deltaY = 1
				}

				if pc.GetY() > aic.TargetY {
					deltaY = -1
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

			//Get if we bumped into something
			entityHit := level.GetSolidEntityAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
			if entityHit == nil {
				tile := level.GetTileAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
				if tile == nil {
				} else if tile.Type != 2 && tile.Type != 4 {
					level.PlaceEntity(pc.GetX()+deltaX, pc.GetY()+deltaY, entity)
				}
			} else {

				//Is it food?
				if entityHit.HasComponent("FoodComponent") && !entityHit.HasComponent("DeadComponent") {
					canEat := false
					if entityHit.HasComponent("HealthComponent") {
						hc := entityHit.GetComponent("HealthComponent").(*component.HealthComponent)
						//Gotta kill the food first.
						if hc.Health > 0 {
							if entityHit.HasComponent("DefensiveAIComponent") {
								daic := entityHit.GetComponent("DefensiveAIComponent").(*component.DefensiveAIComponent)
								daic.Attacked = true
								daic.AttackerX = pc.GetX()
								daic.AttackerY = pc.GetY()
							}
							hc.Health--
						} else {
							canEat = true
						}
					}

					if canEat {
						fc := entityHit.GetComponent("FoodComponent").(*component.FoodComponent)

						if fc.Amount <= 0 {
							entityHit.AddComponent(&component.DeadComponent{})
						} else {
							fc.Amount--
							aic.Energy += 2
						}
					}
				}
			}

			//Change direction we are facing based on movement
			if deltaY > 0 {
				dc.Direction = 1
			}
			if deltaY < 0 {
				dc.Direction = 2
			}
			if deltaX < 0 {
				dc.Direction = 3
			}
			if deltaX > 0 {
				dc.Direction = 0
			}

			//Every thing we just did costed energy.
			aic.Energy--
			if aic.Energy <= 0 {
				entity.AddComponent(&component.DeadComponent{})
			}
		}
	}

	return level
}
