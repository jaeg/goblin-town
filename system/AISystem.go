package system

import (
	"math/rand"

	"goblin-town/world"

	"goblin-town/component"
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

				if entity.HasComponent("GoblinAIComponent") {
					if entity.HasComponent("MyTurnComponent") {
						aic := entity.GetComponent("GoblinAIComponent").(*component.GoblinAIComponent)
						pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
						dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)

						deltaX := 0
						deltaY := 0

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
								fc := entityHit.GetComponent("FoodComponent").(*component.FoodComponent)

								if fc.Amount <= 0 {
									entityHit.AddComponent(&component.DeadComponent{})
								} else {
									fc.Amount--
									aic.Energy += 2
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
