package entity

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jaeg/goblin-town/component"
	"github.com/jaeg/goblin-town/lore"
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
		log.Println(value)
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

	log.Println("Finished", blueprints)
}

func Create(name string, x int, y int) (*Entity, error) {
	blueprint := blueprints[name]
	if blueprint != nil {
		entity := Entity{}
		entity.Blueprint = name
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
			case "InanimateComponent":
				entity.AddComponent(&component.InanimateComponent{})
			case "MassiveComponent":
				entity.AddComponent(&component.MassiveComponent{})
			case "NocturnalComponent":
				entity.AddComponent(&component.NocturnalComponent{})
			case "NeverSleepComponent":
				entity.AddComponent(&component.NeverSleepComponent{})
			case "InventoryComponent":
				inv := &component.InventoryComponent{}
				for _, item := range params {
					inv.AddItem(item)
				}
				entity.AddComponent(inv)
			case "InteractComponent":
				interact := &component.InteractComponent{}
				interact.Message = append(interact.Message, params...)
				log.Println(interact)
				entity.AddComponent(interact)
			case "WanderAIComponent":
				entity.AddComponent(&component.WanderAIComponent{})
			case "GoblinAIComponent":
				entity.AddComponent(&component.GoblinAIComponent{Energy: 100, SightRange: 20, HungerThreshold: 90, State: "wander", SocialThreshold: 4, MateThreshold: 2})
			case "HostileAIComponent":
				r := 5
				if len(params) == 1 {
					r, _ = strconv.Atoi(params[0])
				}

				entity.AddComponent(&component.HostileAIComponent{SightRange: r})
			case "FoodComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.FoodComponent{Amount: amount})
			case "HealthComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.HealthComponent{Health: amount})
			case "DamageComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.DamageComponent{Amount: amount})
			case "PoisonousComponent":
				amount, _ := strconv.Atoi(params[0])
				entity.AddComponent(&component.PoisonousComponent{Duration: amount})
			case "DefensiveAIComponent":
				entity.AddComponent(&component.DefensiveAIComponent{})
			case "DescriptionComponent":
				name := params[0]
				if name == "<GoblinName>" {
					name = lore.RandomGoblinName()
				}
				faction := "none"
				if len(params) == 2 {
					faction = params[1]
				}

				entity.AddComponent(&component.DescriptionComponent{Name: name, Faction: faction})
			}
		}
		return &entity, nil
	}
	return nil, errors.New("no blueprint found")
}
