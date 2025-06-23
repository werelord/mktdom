package main

import "strings"

type categoryType int

const ( // ordered by tab order in finances
	Logistics categoryType = iota
	Appliances
	Mining
	Agriculture
	Scientific
	Military
)

var catList = []categoryType{Logistics, Appliances, Mining, Agriculture, Scientific, Military}

func (c categoryType) String() string {
	return [...]string{"Logistics", "Appliances", "Mining", "Agriculture", "Scientific", "Military"}[c]
}

type productType struct {
	productTypeInternal
	ProductTypeExternal
}
type productTypeInternal struct {
	name     string
	category categoryType
	price    int
}
type ProductTypeExternal struct {
	Supply int
	Demand int
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
			productTypeInternal: productTypeInternal{name: name, category: cat, price: price},
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

type marketVolume struct {
	current int
	total   int
	// share   float32
}

type planet struct {
	Name        string
	Sector      string
	DomPoints   int
	ProductList map[string]productType

	// calculated
	market      marketVolume
	marketByCat map[categoryType]marketVolume
}

// func (p planet) targetMarketShare() float32 {
// 	// target market share to achieve market dominance is essentially the largest market demand for a category
// 	var max = 0
// 	for _, mv := range p.marketByCat {
// 		if mv.total > max {
// 			max = mv.total
// 		}
// 	}

// 	return float32(max) / float32(p.market.total)

// }

var (
	planetMap     = make(map[string]*planet, 0)
	planetDisplay []*planet
	selected      string
)
