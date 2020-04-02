package system

import (
	"fmt"
	"math/rand"

	"goblin-town/component"
	entityFactory "goblin-town/entity"
	"goblin-town/world"
)

func getRandom(low int, high int) int {
	return (rand.Intn((high - low))) + low
}

// PlayerSystem .
func AISystem(planets map[string]*world.Planet) {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			//fmt.Println(t, len(level.Entities))
			for _, entity := range level.Entities {

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
								continue
							}
						}

						deltaX := getRandom(-1, 2)
						deltaY := 0
						if deltaX == 0 {
							deltaY = getRandom(-1, 2)
						}

						if level.GetSolidEntityAt(pc.X+deltaX, pc.Y+deltaY) == nil {
							tile := level.GetTileAt(pc.X+deltaX, pc.Y+deltaY)
							if tile == nil {
							} else if tile.Type != 2 && tile.Type != 4 {
								pc.X += deltaX
								pc.Y += deltaY
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
								continue
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
									daic.AttackerX = pc.X
									daic.AttackerY = pc.Y
								}
							}

							// Point where you attack
							deltaX := 0
							deltaY := 0
							if pc.X < aic.AttackerX {
								deltaX = 1
							}

							if pc.X > aic.AttackerX {
								deltaX = -1
							}

							if pc.Y < aic.AttackerY {
								deltaY = 1
							}

							if pc.Y > aic.AttackerY {
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
								continue
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
							for x := pc.X - aic.SightRange; x < pc.X+aic.SightRange; x++ {
								for y := pc.Y - aic.SightRange; y < pc.Y+aic.SightRange; y++ {
									tile := level.GetTileAt(x, y)
									if tile != nil {
										entityHit := level.GetSolidEntityAt(x, y)
										if entityHit != nil {
											if entityHit.HasComponent("GoblinAIComponent") && !entityHit.HasComponent("DeadComponent") {
												goblinsNearby++
											}
										}
									}
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
								for x := pc.X - 1; x < pc.X+1; x++ {
									for y := pc.Y - 1; y < pc.Y+1; y++ {
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
										planets["hub"].Levels[0].AddEntity(goblin)
									}
								}
							}
						case "findfriends":
							for x := pc.X - aic.SightRange; x < pc.X+aic.SightRange; x++ {
								for y := pc.Y - aic.SightRange; y < pc.Y+aic.SightRange; y++ {
									tile := level.GetTileAt(x, y)
									if tile != nil {
										entityHit := level.GetSolidEntityAt(x, y)
										if entityHit != nil {
											if entityHit.HasComponent("GoblinAIComponent") && !entityHit.HasComponent("DeadComponent") {
												aic.TargetX = x
												aic.TargetY = y
												aic.State = "approach"
												break
											}
										}
									}
								}
							}
						case "approach":
							if pc.X < aic.TargetX {
								deltaX = 1
							}

							if pc.X > aic.TargetX {
								deltaX = -1
							}

							if pc.Y < aic.TargetY {
								deltaY = 1
							}

							if pc.Y > aic.TargetY {
								deltaY = -1
							}

							if deltaX == 0 && deltaY == 0 {
								aic.State = "wander"
							}

						case "search":
							//Scan around for food to the best my vision allows me.
							for x := pc.X - aic.SightRange; x < pc.X+aic.SightRange; x++ {
								for y := pc.Y - aic.SightRange; y < pc.Y+aic.SightRange; y++ {
									tile := level.GetTileAt(x, y)
									if tile != nil {
										entityHit := level.GetSolidEntityAt(x, y)
										if entityHit != nil {
											if entityHit.HasComponent("FoodComponent") && !entityHit.HasComponent("DeadComponent") {
												aic.TargetX = x
												aic.TargetY = y
												aic.State = "approach"
												break
											}
										}
									}
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
						entityHit := level.GetSolidEntityAt(pc.X+deltaX, pc.Y+deltaY)
						if entityHit == nil {
							tile := level.GetTileAt(pc.X+deltaX, pc.Y+deltaY)
							if tile == nil {
							} else if tile.Type != 2 && tile.Type != 4 {
								pc.X += deltaX
								pc.Y += deltaY
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
											daic.AttackerX = pc.X
											daic.AttackerY = pc.Y
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
			}
		}
	}
}
