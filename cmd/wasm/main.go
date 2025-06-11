package main

import (
	"fmt"
	"sort"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

type categoryType int

const (
	Appliances categoryType = iota
	Agriculture
	Logistics
	Military
	Mining
	Scientific
)

type productType int

const (
	transportBot productType = iota
	heavyDutyTransportBot
	transportDrone
	maitenanceDrone
	vacuumBot
	cleaningBot
	lightBot
	screenBot
	snackBot
	personalAssistantRobot
	miningDrone
	minerRobot
	planterDrone
	farmerRobot
	scienceAssistantDrone
	scienceRobot
	combatDrone
	combatRobot
)

func (p productType) String() string {
	return [...]string{"Transport Bot", "Heavy Duty Transport Bot", "Transport Drone", "Maitenance Drone", "Vaccuum Bot", "Cleaning Bot",
		"Light Bot", "Screen Bot", "Snack Bot", "Personal Assistant Robot", "Mining Drone", "Miner Robot", "Planter Drone", "Farmer Robot",
		"Science Assistant Drone", "Science Robot", "Combat Drone", "Combat Robot"}[p]
}

func (p productType) category() categoryType {
	return [...]categoryType{Logistics, Logistics, Logistics, Appliances, Appliances, Appliances,
		Appliances, Appliances, Appliances, Appliances, Mining, Mining, Agriculture, Agriculture,
		Scientific, Scientific, Military, Military}[p]
}

func (p productType) price() int {
	return [...]int{20, 47, 107, 114, 20, 79, 41, 41, 30, 700, 95, 900, 63, 680, 108, 890, 108, 930}[p]
}

type planet struct {
	name         string
	sector       string
	marketVolume int
	marketShare  float32
	domPoints    int
	productList  map[productType]product
}

type product struct {
	supply int
	demand int
}

var planetList = make(map[string]planet, 0)

func main() {
	fmt.Println("go web assembly")
	// js.Global().Set("formatJSON", jsonWrapper())
	js.Global().Set("goLoadData", js.FuncOf(func(this js.Value, args []js.Value) any {
		loadStoredPlanetData()
		displayPlanetList()
		return loadData()
	}))
	js.Global().Set("goOnSectorSelect", js.FuncOf(func(this js.Value, args []js.Value) any {
		return onSectorSelect()
	}))
	js.Global().Set("goAddPlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		return addPlanet()
	}))

	<-make(chan struct{})
}

func loadStoredPlanetData() {
	// todo: load shit
	var p1 = planet{"planet 1", "sector 1", 15415, 14.0, 134, make(map[productType]product, 0)}
	p1.productList[transportBot] = product{4, 6}
	p1.productList[minerRobot] = product{5, 6}
	var p2 = planet{"planet 2", "sector 2", 31337, 57.0, 42, make(map[productType]product, 0)}
	p2.productList[transportBot] = product{4, 6}
	p2.productList[minerRobot] = product{5, 6}
	var p3 = planet{"planet3", "sector 3", 95136, 0.0, 1337, make(map[productType]product, 0)}
	p3.productList[scienceRobot] = product{4, 6}
	p3.productList[scienceAssistantDrone] = product{5, 6}

	planetList[p1.name] = p1
	planetList[p2.name] = p2
	planetList[p3.name] = p3
}

// func prettyJson(input string) (string, error) {
// 	var raw any
// 	if err := json.Unmarshal([]byte(input), &raw); err != nil {
// 		return "", err
// 	}
// 	pretty, err := json.MarshalIndent(raw, "", "  ")
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(pretty), nil
// }

// func jsonWrapper() js.Func {
// 	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
// 		if len(args) != 1 {
// 			return sendErr("Invalid no of arguments passed")
// 		}

// 		jsDoc := js.Global().Get("document")
// 		if !jsDoc.Truthy() {
// 			return sendErr("unable to get document object")
// 		}
// 		jsonOutputTextArea := jsDoc.Call("getElementById", "jsonoutput")
// 		if !jsonOutputTextArea.Truthy() {
// 			return sendErr("unable to get output text area")
// 		}
// 		inputJSON := args[0].String()
// 		fmt.Printf("input %s\n", inputJSON)
// 		pretty, err := prettyJson(inputJSON)
// 		if err != nil {
// 			errStr := fmt.Sprintf("unable to convert to json %s\n", err)
// 			return sendErr(errStr)
// 		}
// 		jsonOutputTextArea.Set("value", pretty)
// 		js.Global().Get("localStorage").Call("setItem", "foo", pretty)
// 		return nil
// 	})
// 	return jsonFunc
// }

func loadData() any {
	// todo: populate
	fmt.Println("in loadData")
	// loadSectorList()
	return sendErr("testing")
	// return nil
}

func sendErr(val string) map[string]any {
	res := map[string]any{
		"error": val,
	}
	return res
}

func displayPlanetList() any {
	doc := dom.GetWindow().Document()
	el := doc.GetElementByID("planetList")
	fmt.Printf("full planet list: %v\n", el.Class())
	for _, child := range el.ChildNodes() {
		fmt.Printf("child %v - %v\n", child.NodeName(), child.NodeValue())
		el.RemoveChild(child)
	}

	planetkeys := make([]string, 0, len(planetList))
	for k := range planetList {
		planetkeys = append(planetkeys, k)
	}
	sort.Strings(planetkeys)

	for _, pk := range planetkeys {
		planet := planetList[pk]

		fullPlanetDiv := doc.CreateElement("div")
		fullPlanetDiv.SetAttribute("class", "fullPlanetDetails")

		planetWrapperDiv := doc.CreateElement("div")
		planetWrapperDiv.SetAttribute("onClick", "switchSelected(this)")
		fullPlanetDiv.AppendChild(planetWrapperDiv)

		nameDiv := doc.CreateElement("div")
		nameDiv.SetAttribute("class", "planetName")
		nameDiv.SetInnerHTML(planet.name)
		planetWrapperDiv.AppendChild(nameDiv)

		detailsDiv := doc.CreateElement("div")
		detailsDiv.SetAttribute("class", "planetDetails")
		planetWrapperDiv.AppendChild(detailsDiv)

		sectorDiv := doc.CreateElement("div")
		sectorDiv.SetInnerHTML(fmt.Sprintf("sector %v", planet.sector))
		marketVolDiv := doc.CreateElement("div")
		marketVolDiv.SetInnerHTML(fmt.Sprintf("market cap: %v", planet.marketVolume))
		marketShareDiv := doc.CreateElement("div")
		marketShareDiv.SetInnerHTML(fmt.Sprintf("market share: %.2f%%", planet.marketShare))
		pointsDiv := doc.CreateElement("div")
		pointsDiv.SetInnerHTML(fmt.Sprintf("points: %v", planet.domPoints))
		detailsDiv.AppendChild(sectorDiv)
		detailsDiv.AppendChild(marketVolDiv)
		detailsDiv.AppendChild(marketShareDiv)
		detailsDiv.AppendChild(pointsDiv)

		marketDiv := doc.CreateElement("div")
		marketDiv.SetAttribute("class", "planetMarket")
		fullPlanetDiv.AppendChild(marketDiv)

		for t, prod := range planet.productList {
			prodDiv := doc.CreateElement("div")
			prodDiv.SetInnerHTML(fmt.Sprintf("%v: %v/%v", t.String(), prod.supply, prod.demand))
			marketDiv.AppendChild(prodDiv)
		}

		el.AppendChild(fullPlanetDiv)
	}
	return nil
}

func onSectorSelect() any {
	return sendErr("not yet implemented")
}
func addPlanet() any {
	return sendErr("not yet implemented")
}
func (p planet) calcMarketCap() int {
	fmt.Println("not yet implemented")
	return 0
}

func (p planet) calcMarketShare() float32 {
	return 0
}

// func getSectorList() []string {
// 	var list = []string{"sector 1", "sector 2", "foobar"}
// 	return list
// }
