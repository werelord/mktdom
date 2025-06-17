package main

import "strings"

type categoryType int

const (
	Appliances categoryType = iota
	Agriculture
	Logistics
	Military
	Mining
	Scientific
)

// const (
// 	transportBot productType = iota
// 	heavyDutyTransportBot
// 	transportDrone
// 	maitenanceDrone
// 	vacuumBot
// 	cleaningBot
// 	lightBot
// 	screenBot
// 	snackBot
// 	personalAssistantRobot
// 	miningDrone
// 	minerRobot
// 	planterDrone
// 	farmerRobot
// 	scienceAssistantDrone
// 	scienceRobot
// 	combatDrone
// 	combatRobot
// )

type productType struct {
	productTypeExport
	productTypeInternal
}

type productTypeExport struct {
	Supply int
	Demand int
}

type productTypeInternal struct {
	name     string
	category categoryType
	price    int
}

func (p productType) ID() string {
	return toId(p.name)
}

func toId(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "_")
}

var baseProductMap map[string]productType

func init() {

	baseProductMap = make(map[string]productType, 18)

	var f = func(name string, cat categoryType, price int) productType {
		var prod = productType{
			productTypeExport{},
			productTypeInternal{name: name, category: cat, price: price},
		}
		return prod
	}

	var prodList = []productType{
		f("Transport Bot", Logistics, 20),
		f("Heavy Duty Transport Bot", Logistics, 47),
		f("Transport Drone", Logistics, 107),

		f("Maitenance Drone", Appliances, 114),
		f("Vaccuum Bot", Appliances, 20),
		f("Cleaning Bot", Appliances, 79),
		f("Light Bot", Appliances, 41),
		f("Screen Bot", Appliances, 41),
		f("Snack Bot", Appliances, 30),
		f("Personal Assistant Robot", Appliances, 700),

		f("Mining Drone", Mining, 95),
		f("Miner Robot", Mining, 900),

		f("Planter Drone", Agriculture, 63),
		f("Farmer Robot", Agriculture, 680),

		f("Science Assistant Drone", Scientific, 108),
		f("Science Robot", Scientific, 890),
		f("Combat Drone", Military, 108),
		f("Combat Robot", Military, 930),
	}

	for _, prod := range prodList {
		baseProductMap[prod.ID()] = prod
	}
}

// convenience for form building.. builds two columns, category downwards
var productList = []string{
	toId("Transport Bot"),
	toId("Maitenance Drone"),
	toId("Heavy Duty Transport Bot"),
	toId("Vaccuum Bot"),
	toId("Transport Drone"),
	toId("Cleaning Bot"),
	toId("Mining Drone"),
	toId("Light Bot"),
	toId("Miner Robot"),
	toId("Screen Bot"),
	toId("Planter Drone"),
	toId("Snack Bot"),
	toId("Farmer Robot"),
	toId("Personal Assistant Robot"),
	toId("Science Assistant Drone"),
	toId("Combat Drone"),
	toId("Science Robot"),
	toId("Combat Robot"),
}

type planet struct {
	name         string
	sector       string
	marketVolume int
	marketShare  float32
	domPoints    int
	productList  map[string]productType
}

var planetList = make(map[string]planet, 0)
