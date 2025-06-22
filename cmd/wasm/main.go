package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

func main() {
	fmt.Println("go web assembly")
	// js.Global().Set("formatJSON", jsonWrapper())
	js.Global().Set("goLoadData", js.FuncOf(func(this js.Value, args []js.Value) any {
		return loadData()
	}))
	js.Global().Set("goOnSelected", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return sendErr("error in setting selected (wrong # of args)")
		}
		return onSelected(args[0].String())
	}))
	js.Global().Set("goAddPlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		return addPlanet()
	}))
	js.Global().Set("goGenPlanetForm", js.FuncOf(func(this js.Value, args []js.Value) any {
		return genPlanetForm()
	}))

	<-make(chan struct{})
}
func sendErr(val string) map[string]any {
	res := map[string]any{
		"error": val,
	}
	return res
}

func sendToast(t string) map[string]any {
	res := map[string]any{
		"toast": t,
	}
	return res
}

func loadData() any {
	// todo: populate

	loadStoredPlanetData()
	// fmt.Println("skipping loading display")
	loadPlanetDisplay()
	genPlanetForm()

	// return sendErr("testing")
	return nil
}

func loadStoredPlanetData() {

	if obj := js.Global().Get("localStorage").Call("getItem", "planetMap"); obj.IsNull() {
		// todo: nothing in storage; load example data??
	} else {
		// fmt.Printf("load successful, json: %v\n", obj.String())
		if err := json.Unmarshal([]byte(obj.String()), &planetMap); err != nil {
			fmt.Printf("error unmarshalling: %v\n", err)
		} else {
			// fmt.Printf("unmarshal successful: %+v\n", planetMap)
			// set the base internal (non-exported) values
			for _, planet := range planetMap {
				for pname, prod := range planet.ProductList {
					prod.productTypeInternal = baseProductMap[pname].productTypeInternal
					planet.ProductList[pname] = prod
				}
				planet.calcMarketVol()
			}
		}
	}
}

func savePlanetData() error {
	fmt.Println("in savePlanetData")
	str, err := json.Marshal(planetMap)
	if err != nil {
		fmt.Printf("error :%v\n", err)
		return err
	} else {
		// fmt.Printf("json: %v\n", string(str))
		js.Global().Get("localStorage").Call("setItem", "planetMap", string(str))
		fmt.Println("savePlanetData successful")
		return nil
	}
}

func loadPlanetDisplay() any {

	doc := dom.GetWindow().Document()
	planetListDiv := doc.GetElementByID("planetList")
	// fmt.Printf("full planet list: %v\n", el.Class())
	for _, child := range planetListDiv.ChildNodes() {
		// fmt.Printf("child %v - %v\n", child.NodeName(), child.NodeValue())
		planetListDiv.RemoveChild(child)
	}

	planetDisplay = make([]*planet, 0, len(planetMap))
	for _, p := range planetMap {
		// fmt.Println("formap:" + p.Name)
		planetDisplay = append(planetDisplay, p)
	}
	sort.Slice(planetDisplay, func(i, j int) bool {
		return planetDisplay[i].Name < planetDisplay[j].Name
	})

	for _, planet := range planetDisplay {
		fullPlanetDiv := generatePlanetDisplay(*planet)
		planetListDiv.AppendChild(fullPlanetDiv)
	}
	return nil
}

func generatePlanetDisplay(p planet) dom.Element {

	doc := dom.GetWindow().Document()

	// full planet wrapper
	fullPlanetDiv := doc.CreateElement("div")
	// fullPlanetDiv.SetID(p.Name)
	fullPlanetDiv.Class().SetString("fullPlanetDetails")
	fullPlanetDiv.SetAttribute("onClick", fmt.Sprintf("switchSelected('%v')", p.Name))

	planetWrapperDiv := doc.CreateElement("div")
	planetWrapperDiv.SetID(p.Name)
	planetWrapperDiv.Class().SetString("planetInfo")
	// planetWrapperDiv.SetAttribute("onClick", "switchSelected(this)")
	fullPlanetDiv.AppendChild(planetWrapperDiv)

	nameDiv := doc.CreateElement("div")
	nameDiv.Class().SetString("planetName")
	nameDiv.SetInnerHTML(p.Name)
	planetWrapperDiv.AppendChild(nameDiv)

	detailsDiv := doc.CreateElement("div")
	detailsDiv.Class().SetString("planetDetails")
	planetWrapperDiv.AppendChild(detailsDiv)

	sectorDiv := doc.CreateElement("div")
	sectorDiv.SetInnerHTML(fmt.Sprintf("sector %v", p.Sector))
	marketVolDiv := doc.CreateElement("div")
	marketVolDiv.SetInnerHTML(fmt.Sprintf("market cap: %v", p.market.total))
	marketShareDiv := doc.CreateElement("div")
	marketShareDiv.SetInnerHTML(fmt.Sprintf("market share: %.2f%%", (float32(p.market.current) / float32(p.market.total) * 100)))
	pointsDiv := doc.CreateElement("div")
	pointsDiv.SetInnerHTML(fmt.Sprintf("points: %v", p.DomPoints))
	detailsDiv.AppendChild(sectorDiv)
	detailsDiv.AppendChild(marketVolDiv)
	detailsDiv.AppendChild(pointsDiv)
	detailsDiv.AppendChild(marketShareDiv)

	marketDiv := doc.CreateElement("div")
	marketDiv.Class().SetString("planetMarket")
	fullPlanetDiv.AppendChild(marketDiv)

	for _, prod := range p.ProductList {
		prodDiv := doc.CreateElement("div")
		prodDiv.SetInnerHTML(fmt.Sprintf("%v: %v/%v", prod.name, prod.Supply, prod.Demand))
		marketDiv.AppendChild(prodDiv)
	}

	categoryDiv := doc.CreateElement("div")
	categoryDiv.Class().SetString("planetCategoryMarket")
	// var (
	// 	max   int
	// 	maxEl dom.Element
	// )

	// for cat, catMarket := range p.marketByCat {
	for _, cat := range catList {
		if catMarket, exists := p.marketByCat[cat]; exists {

			catMarketDiv := doc.CreateElement("div")
			catMarketDiv.SetID(cat.String())

			var (
				cur = float64(catMarket.current) / float64(p.market.total) * 100
				opp = (float64(catMarket.total) - float64(catMarket.current)) / float64(p.market.total) * 100
			)

			cur = math.Round(cur*100) / 100
			opp = math.Round(opp*100) / 100

			catMarketDiv.SetInnerHTML(fmt.Sprintf("%v: %.1f / %.1f%%", cat.String(), cur, opp))

			// for highlighting
			// if catMarket.total > max {
			// 	max = catMarket.total
			// 	maxEl = catMarketDiv
			// }
			categoryDiv.AppendChild(catMarketDiv)
		}
	}
	// maxEl.Class().Add("highlight")
	fullPlanetDiv.AppendChild(categoryDiv)

	return fullPlanetDiv
}

func onSelected(newSel string) any {
	// fmt.Printf("onSelected, id=%v\n", newSel)

	var (
		prevSel = selected
		doc     = dom.GetWindow().Document()
	)

	if prevSel != "" {
		var oldDiv = doc.GetElementByID(prevSel)
		oldDiv.Class().Remove("selected")
	}

	if newSel == prevSel {
		// same object, its just a deselect
		selected = ""
	} else {
		var newDiv = doc.GetElementByID(newSel)
		newDiv.Class().Add("selected")
		selected = newSel
	}

	// todo: other stuff?

	return nil
}

// func (p *planet) calcMarketShare() float32 {
// 	return 0
// }

// func getSectorList() []string {
// 	var list = []string{"sector 1", "sector 2", "foobar"}
// 	return list
// }
