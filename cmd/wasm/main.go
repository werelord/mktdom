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
	js.Global().Set("goOnAddPlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		return onAddPlanet()
	}))
	js.Global().Set("goShowAddPlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return sendErr("error in showing Add Planet form (wrong # of args)")
		}
		return showAddPlanet(args[0].String())
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

	planetListDiv.AppendChild(genPlanetHeader(doc))

	planetDisplay = make([]*planet, 0, len(planetMap))
	for _, p := range planetMap {
		// fmt.Println("formap:" + p.Name)
		planetDisplay = append(planetDisplay, p)
	}

	// sort by name
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
	fullPlanetDiv.SetID(p.Name)
	fullPlanetDiv.Class().SetString("fullPlanetDetails")
	fullPlanetDiv.SetAttribute("onClick", fmt.Sprintf("switchSelected('%v')", p.Name))

	fullPlanetDiv.AppendChild(genPlanetInfo(doc, p))

	fullPlanetDiv.AppendChild(genPlanetMarket(doc, p))

	fullPlanetDiv.AppendChild(genPlanetCatMarket(doc, p))

	return fullPlanetDiv
}

func genPlanetHeader(doc dom.Document) dom.Element {

	headerdiv := doc.CreateElement("div")
	headerdiv.Class().SetString("fullPlanetDetails")

	info := doc.CreateElement("div")
	info.Class().SetString("planetInfoHeader")
	info.SetInnerHTML("Planet Information")

	market := doc.CreateElement("div")
	market.Class().SetString("planetMarketHeader")
	market.SetInnerHTML("Products (my amt / total amt)")

	category := doc.CreateElement("div")
	category.Class().SetString("planetCategoryMarketHeader")
	category.SetInnerHTML("Category (my share / opp share %)")

	headerdiv.AppendChild(info)
	headerdiv.AppendChild(market)
	headerdiv.AppendChild(category)

	return headerdiv
}

func genPlanetInfo(doc dom.Document, p planet) dom.Element {

	planetWrapperDiv := doc.CreateElement("div")
	// planetWrapperDiv.SetID(p.Name)
	planetWrapperDiv.Class().SetString("planetInfo")
	// planetWrapperDiv.SetAttribute("onClick", "switchSelected(this)")

	nameRow := doc.CreateElement("div")
	nameRow.Class().SetString("planetNameRow")
	planetWrapperDiv.AppendChild(nameRow)

	nameDiv := doc.CreateElement("div")
	nameDiv.Class().SetString("planetName")
	nameDiv.SetInnerHTML(p.Name)
	nameRow.AppendChild(nameDiv)

	editimg := doc.CreateElement("img")
	editimg.Class().SetString("editImage")
	editimg.SetAttribute("src", "img/pencil-square-o.svg")
	editimg.SetAttribute("onclick", fmt.Sprintf("editPlanet('%v', event);", p.Name))
	nameRow.AppendChild(editimg)

	detailsDiv := doc.CreateElement("div")
	detailsDiv.Class().SetString("planetDetails")
	planetWrapperDiv.AppendChild(detailsDiv)

	sectorDiv := doc.CreateElement("div")
	sectorDiv.SetInnerHTML(fmt.Sprintf("sector %v", p.Sector))
	marketVolDiv := doc.CreateElement("div")
	marketVolDiv.SetInnerHTML(fmt.Sprintf("market cap: %v", p.market.total))
	marketShareDiv := doc.CreateElement("div")
	marketShareDiv.SetID("market_share")
	marketShareDiv.SetInnerHTML(fmt.Sprintf("market share: %.2f%%", (float32(p.market.current) / float32(p.market.total) * 100)))
	pointsDiv := doc.CreateElement("div")
	pointsDiv.SetInnerHTML(fmt.Sprintf("points: %v", p.DomPoints))
	detailsDiv.AppendChild(sectorDiv)
	detailsDiv.AppendChild(marketVolDiv)
	detailsDiv.AppendChild(pointsDiv)
	detailsDiv.AppendChild(marketShareDiv)

	return planetWrapperDiv
}

func genPlanetMarket(doc dom.Document, p planet) dom.Element {

	marketDiv := doc.CreateElement("div")
	marketDiv.Class().SetString("planetMarket")

	// want to order this correctly.. similar to addPlanetForm
	// this won't be exact, since not all items will be in the list
	// but at least it will try to keep it consistent

	for _, prodId := range productList {
		if prod, exist := p.ProductList[prodId]; exist {
			prodDiv := doc.CreateElement("div")
			prodDiv.SetID(prod.ID())
			prodDiv.SetInnerHTML(fmt.Sprintf("%v: %v/%v", prod.name, prod.Supply, prod.Demand))
			marketDiv.AppendChild(prodDiv)
		}
	}

	return marketDiv

}

func genPlanetCatMarket(doc dom.Document, p planet) dom.Element {
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

	return categoryDiv

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
