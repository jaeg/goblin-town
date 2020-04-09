package entity

import (
	"bufio"
	"errors"
	"fmt"
	"goblin-town/component"
	"os"
	"strconv"
	"strings"
)

var blueprints = make(map[string][]string)

// FactoryLoad Loads the blueprints for the factory to construct entities
func FactoryLoad(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	entityName := ""
	for scanner.Scan() {
		value := scanner.Text()
		fmt.Println(value)
		if value == "" {
			entityName = ""
			continue
		}
		if entityName == "" {
			entityName = value
			continue
		} else {
			blueprints[entityName] = append(blueprints[entityName], value)
		}
	}

	fmt.Println("Finished", blueprints)
}

func Create(name string, x int, y int) (*Entity, error) {
	blueprint := blueprints[name]
	if blueprint != nil {
		entity := Entity{}

		pc := &component.PositionComponent{}
		pc.SetPosition(x, y)
		entity.AddComponent(pc)

		entity.AddComponent(&component.DirectionComponent{Direction: 0})
		for _, value := range blueprint {
			c := strings.Split(value, ":")
			params := strings.Split(c[1], ",")
			switch c[0] {
			case "AppearanceComponent":
				sx, _ := strconv.Atoi(params[0])
				sy, _ := strconv.Atoi(params[1])
				r, _ := strconv.Atoi(params[2])
				g, _ := strconv.Atoi(params[3])
				b, _ := strconv.Atoi(params[4])
				entity.AddComponent(&component.AppearanceComponent{SpriteX: int32(sx), SpriteY: int32(sy), R: uint8(r), G: uint8(g), B: uint8(b)})

			case "InitiativeComponent":
				dv, _ := strconv.Atoi(params[0])
				ticks, _ := strconv.Atoi(params[1])
				entity.AddComponent(&component.InitiativeComponent{DefaultValue: dv, Ticks: ticks})
			case "SolidComponent":
				entity.AddComponent(&component.SolidComponent{})
			case "InventoryComponent":
				inv := &component.InventoryComponent{}
				for _, item := range params {
					inv.AddItem(item)
				}
				entity.AddComponent(inv)
			case "InteractComponent":
				interact := &component.InteractComponent{}
				for _, item := range params {
					interact.Message = append(interact.Message, item)
				}
				fmt.Println(interact)
				entity.AddComponent(interact)
			case "WanderAIComponent":
				entity.AddComponent(&component.WanderAIComponent{})
			case "GoblinAIComponent":
				entity.AddComponent(&component.GoblinAIComponent{Energy: 100, SightRange: 20, HungerThreshold: 90, State: "wander", SocialThreshold: 4, MateThreshold: 2})
			case "FoodComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.FoodComponent{Amount: amount})
			case "HealthComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.HealthComponent{Health: amount})
			case "DefensiveAIComponent":
				entity.AddComponent(&component.DefensiveAIComponent{})
			}
		}
		return &entity, nil
	}
	return nil, errors.New("No blueprint found")
}
