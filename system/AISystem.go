package system

import (
	"fmt"
	"math/rand"

	"goblin-town/component"
	"goblin-town/entity"
	"goblin-town/world"
)

func getRandom(low int, high int) int {
	return (rand.Intn((high - low))) + low
}

type AISystem struct {
}

// PlayerSystem .
func (s AISystem) Update(level *world.Level, entity *entity.Entity) *world.Level {
	if entity.HasComponent("WanderAIComponent") {
		if entity.HasComponent("MyTurnComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)

			if entity.HasComponent("HealthComponent") {
				hc := entity.GetComponent("HealthComponent").(*component.HealthComponent)
				if hc.Health <= 0 {
					if entity.HasComponent("FoodComponent") {
						entity.RemoveComponent("WanderAIComponent")
					} else {
						entity.AddComponent(&component.DeadComponent{})
					}
					return level
				}
			}

			deltaX := getRandom(-1, 2)
			deltaY := 0
			if deltaX == 0 {
				deltaY = getRandom(-1, 2)
			}

			if level.GetSolidEntityAt(pc.GetX()+deltaX, pc.GetY()+deltaY) == nil {
				tile := level.GetTileAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
				if tile == nil {
				} else if tile.Type != 2 && tile.Type != 4 {
					level.PlaceEntity(pc.GetX()+deltaX, pc.GetY()+deltaY, entity)

				}
			}
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
		}
	}

	if entity.HasComponent("HostileAIComponent") {
		if entity.HasComponent("MyTurnComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
			hc := entity.GetComponent("HostileAIComponent").(*component.HostileAIComponent)
			deltaX := 0
			deltaY := 0

			// Handle being dead
			if entity.HasComponent("HealthComponent") {
				hc := entity.GetComponent("HealthComponent").(*component.HealthComponent)
				if hc.Health <= 0 {
					if entity.HasComponent("FoodComponent") {
						entity.RemoveComponent("HostileAIComponent")
					} else {
						entity.AddComponent(&component.DeadComponent{})
					}
					return level
				}
			}
			//Scan around for food to the best my vision allows me.
			nearby := level.GetEntitiesAround(pc.GetX(), pc.GetY(), hc.SightRange, hc.SightRange)
			hunting := false
			for e := range nearby {
				if nearby[e] != entity {
					if (nearby[e].HasComponent("FoodComponent") || nearby[e].HasComponent("GoblinAIComponent")) && !nearby[e].HasComponent("DeadComponent") {
						foodPC := nearby[e].GetComponent("PositionComponent").(*component.PositionComponent)
						hc.TargetX = foodPC.GetX()
						hc.TargetY = foodPC.GetY()
						hunting = true
						break
					}
				}
			}

			if hunting {
				if pc.GetX() < hc.TargetX {
					deltaX = 1
				}

				if pc.GetX() > hc.TargetX {
					deltaX = -1
				}

				if pc.GetY() < hc.TargetY {
					deltaY = 1
				}

				if pc.GetY() > hc.TargetY {
					deltaY = -1
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
			entityHit := level.GetSolidEntityAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
			if entityHit == nil {
				tile := level.GetTileAt(pc.GetX()+deltaX, pc.GetY()+deltaY)
				if tile == nil {
				} else if tile.Type != 2 && tile.Type != 4 {
					level.PlaceEntity(pc.GetX()+deltaX, pc.GetY()+deltaY, entity)

				}
			} else {
				if entityHit != entity {
					//Attack it
					if entityHit.HasComponent("HealthComponent") {
						ehc := entityHit.GetComponent("HealthComponent").(*component.HealthComponent)
						if ehc.Health > 0 {
							ehc.Health--
						}

						if entityHit.HasComponent("FoodComponent") {
							fc := entityHit.GetComponent("FoodComponent").(*component.FoodComponent)

							if fc.Amount <= 0 {
								entityHit.AddComponent(&component.DeadComponent{})
							} else {
								fc.Amount--
							}
						}
					}

				}

				// Trigger their defenses
				if entityHit.HasComponent("DefensiveAIComponent") {
					daic := entityHit.GetComponent("DefensiveAIComponent").(*component.DefensiveAIComponent)
					daic.Attacked = true
					daic.AttackerX = pc.GetX()
					daic.AttackerY = pc.GetY()
				}
			}
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
		}
	}

	if entity.HasComponent("DefensiveAIComponent") {
		if entity.HasComponent("MyTurnComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
			aic := entity.GetComponent("DefensiveAIComponent").(*component.DefensiveAIComponent)

			// Handle being dead
			if entity.HasComponent("HealthComponent") {
				hc := entity.GetComponent("HealthComponent").(*component.HealthComponent)
				if hc.Health <= 0 {
					if entity.HasComponent("FoodComponent") {
						entity.RemoveComponent("DefensiveAIComponent")
					} else {
						entity.AddComponent(&component.DeadComponent{})
					}
					return level
				}
			}

			if aic.Attacked {
				entityHit := level.GetSolidEntityAt(aic.AttackerX, aic.AttackerY)

				if entityHit == nil {
					// No attacker there.
					aic.Attacked = false
				} else {
					// Hit the attacker back.
					if entityHit.HasComponent("HealthComponent") {
						ehc := entityHit.GetComponent("HealthComponent").(*component.HealthComponent)
						if ehc.Health > 0 {
							ehc.Health--
							fmt.Println("Health left", ehc.Health)
							aic.Attacked = false
						}
					}

					// Trigger their defenses
					if entityHit.HasComponent("DefensiveAIComponent") {
						daic := entityHit.GetComponent("DefensiveAIComponent").(*component.DefensiveAIComponent)
						daic.Attacked = true
						daic.AttackerX = pc.GetX()
						daic.AttackerY = pc.GetY()
					}
				}

				// Point where you attack
				deltaX := 0
				deltaY := 0
				if pc.GetX() < aic.AttackerX {
					deltaX = 1
				}

				if pc.GetX() > aic.AttackerX {
					deltaX = -1
				}

				if pc.GetY() < aic.AttackerY {
					deltaY = 1
				}

				if pc.GetY() > aic.AttackerY {
					deltaY = -1
				}

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
			}
		}
	}

	return level
}
