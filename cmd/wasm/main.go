package main

import (
	"fmt"
	"sort"
	"strings"
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

func main() {
	fmt.Println("go web assembly")
	// js.Global().Set("formatJSON", jsonWrapper())
	js.Global().Set("goLoadData", js.FuncOf(func(this js.Value, args []js.Value) any {
		loadStoredPlanetData()
		displayPlanetList()
		genPlanetForm()
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
	tb := baseProductMap["transport_bot"]
	tb.productTypeExport = productTypeExport{4, 6}
	mr := baseProductMap["miner_robot"]
	mr.productTypeExport = productTypeExport{5, 6}

	var p1 = planet{"planet 331", "sector 1", 15415, 14.0, 134, make(map[string]productType, 0)}
	p1.productList[tb.ID()] = tb
	p1.productList[mr.ID()] = mr

	var p2 = planet{"planet 2", "sector 2", 31337, 57.0, 42, make(map[string]productType, 0)}
	p2.productList[tb.ID()] = tb
	p2.productList[mr.ID()] = mr

	var p3 = planet{"planet3", "sector 3", 95136, 0.0, 1337, make(map[string]productType, 0)}
	p3.productList[mr.ID()] = mr
	p3.productList[tb.ID()] = tb
	var p4 = planet{"planet 4", "sector 1", 15415, 14.0, 134, make(map[string]productType, 0)}
	p4.productList[tb.ID()] = tb
	p4.productList[mr.ID()] = mr
	var p5 = planet{"planet 5", "sector 2", 31337, 57.0, 42, make(map[string]productType, 0)}
	p5.productList[tb.ID()] = tb
	p5.productList[mr.ID()] = mr
	var p6 = planet{"planet6", "sector 3", 95136, 0.0, 1337, make(map[string]productType, 0)}
	p6.productList[tb.ID()] = tb
	p6.productList[mr.ID()] = mr

	planetList[p1.name] = p1
	planetList[p2.name] = p2
	planetList[p3.name] = p1
	planetList[p4.name] = p4
	planetList[p5.name] = p5
	planetList[p6.name] = p6

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
	// fmt.Printf("full planet list: %v\n", el.Class())
	for _, child := range el.ChildNodes() {
		// fmt.Printf("child %v - %v\n", child.NodeName(), child.NodeValue())
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
		detailsDiv.AppendChild(pointsDiv)
		detailsDiv.AppendChild(marketShareDiv)

		marketDiv := doc.CreateElement("div")
		marketDiv.SetAttribute("class", "planetMarket")
		fullPlanetDiv.AppendChild(marketDiv)

		for _, prod := range planet.productList {
			prodDiv := doc.CreateElement("div")
			prodDiv.SetInnerHTML(fmt.Sprintf("%v: %v/%v", prod.name, prod.Supply, prod.Demand))
			marketDiv.AppendChild(prodDiv)
		}

		el.AppendChild(fullPlanetDiv)
	}
	return nil
}

func genPlanetForm() any {
	// return nil
	doc := dom.GetWindow().Document()
	table := doc.GetElementByID("addPlanetTable")

	var genProdTd = func(r dom.Element, p productType) {
		// var tdlist = make([]dom.Element, 0)

		label := doc.CreateElement("td")
		label.SetAttribute("class", "productLabel")
		// label.SetAttribute("id", fmt.Sprintf("%v", spToUl(p.ID())))
		// fmt.Println("in genPlanetForm, " + p.name)
		label.SetInnerHTML(p.name)
		r.AppendChild(label)

		mbntd := doc.CreateElement("td")
		minus := doc.CreateElement("button")
		minus.SetAttribute("class", "productButton")
		minus.SetAttribute("onsubmit", "return false;")
		minus.SetInnerHTML("-")
		mbntd.AppendChild(minus)
		r.AppendChild(mbntd)

		amt := doc.CreateElement("td")
		amt.SetAttribute("class", "productAmount")
		amt.SetAttribute("id", fmt.Sprintf("%v_amt", p.ID()))
		amt.SetInnerHTML(fmt.Sprintf("%v", 0))
		r.AppendChild(amt)

		pbntd := doc.CreateElement("td")
		plus := doc.CreateElement("button")
		plus.SetAttribute("class", "productButton")
		plus.SetAttribute("onsubmit", "return false;")
		plus.SetInnerHTML("+")
		pbntd.AppendChild(plus)
		r.AppendChild(pbntd)

	}

	for i := 0; i < len(productList); i = i + 2 {
		var p1, p2 = baseProductMap[productList[i]], baseProductMap[productList[i+1]]

		row := doc.CreateElement("tr")

		genProdTd(row, p1)
		genProdTd(row, p2)

		// plus := doc.CreateElement("button")

		table.AppendChild(row)
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

func spToUl(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}
