package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return onAddPlanet(false)
	}))
	js.Global().Set("goOnSavePlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		return onAddPlanet(true)
	}))
	js.Global().Set("goOnDeletePlanet", js.FuncOf(func(this js.Value, args []js.Value) any {
		return onDeletePlanet()
	}))
	js.Global().Set("goGenAddPlanetForm", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return sendErr("error in showing Add Planet form (wrong # of args)")
		}
		return genAddPlanetForm(args[0].String())
	}))
	js.Global().Set("goOnSubSupply", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return sendErr("error in changing product supply (wrong # of args)")
		}
		return onChangeSupply(-48, args[0].String())
	}))
	js.Global().Set("goOnAddSupply", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return sendErr("error in changing product supply (wrong # of args)")
		}
		return onChangeSupply(48, args[0].String())
	}))

	<-make(chan struct{})
}
func sendErr(f string, val ...any) map[string]any {
	res := map[string]any{
		"error": fmt.Sprintf(f, val...),
	}
	return res
}

func sendToast(t string, val ...any) map[string]any {
	res := map[string]any{
		"toast": fmt.Sprintf(t, val...),
	}
	return res
}

func loadData() any {
	if err := loadStoredPlanetData(); err != nil {
		return sendErr("error: %v", err)
	}

	loadPlanetListDisplay()
	return nil
}

func loadStoredPlanetData() error {

	var (
		localstorage = js.Global().Get("localStorage")
		joinErr      error
	)
	// fmt.Printf("localstorage length = %v\n", localstorage.Length())
	planetMap = make(map[string]*planetType, localstorage.Length())
	for i := 0; i < localstorage.Length(); i++ {
		var (
			key    = localstorage.Call("key", i).String()
			planet planetType
		)

		if obj := localstorage.Call("getItem", key); obj.IsNull() {
			joinErr = errors.Join(joinErr, fmt.Errorf("load error: object is null: '%v'", key))
		} else if err := json.Unmarshal([]byte(obj.String()), &planet); err != nil {
			joinErr = errors.Join(joinErr, fmt.Errorf("error unmarshalling: %v\n", err))
		} else {
			for name, prod := range planet.ProductList {
				// make sure base are inserted back
				prod.productTypeInternal = baseProductMap[name].productTypeInternal
				planet.ProductList[name] = prod
			}
			planet.calcMarketVol()
			// fmt.Printf("loaded planet '%v'\n", planet)
			planetMap[planet.Name] = &planet
		}
	}

	// temporary
	// for _, planet := range planetMap {
	// 	deletePlanetData(planet.Name)

	// 	planet.Name = strings.Title(planet.Name)
	// 	planet.Sector = strings.ToUpper(planet.Sector)

	// 	savePlanetData(*planet)
	// }
	return joinErr
}

// func saveAllPlanetData() error {
// 	fmt.Println("in savePlanetData")

// 	var localstorage = js.Global().Get("localStorage")

// 	for _, planet := range planetMap {
// 		if str, err := json.Marshal(*planet); err != nil {
// 			return err
// 		} else {
// 			// fmt.Printf("%v\n", string(str))
// 			localstorage.Call("setItem", planet.Name, string(str))
// 		}
// 	}

// 	fmt.Println("savePlanetData successful")
// 	return nil
// }

func savePlanetData(planet planetType) error {
	fmt.Println("in savePlanetData")
	var localstorage = js.Global().Get("localStorage")
	if str, err := json.Marshal(planet); err != nil {
		return err
	} else {
		// fmt.Printf("%v\n", string(str))
		localstorage.Call("setItem", planet.Name, string(str))
	}

	return nil
}

func deletePlanetData(name string) {
	fmt.Println("in delete planet")

	var localstorage = js.Global().Get("localStorage")
	localstorage.Call("removeItem", name)

}

func loadPlanetListDisplay() any {

	doc := dom.GetWindow().Document()
	planetListDiv := doc.GetElementByID("planetList")
	// fmt.Printf("full planet list: %v\n", el.Class())
	for _, child := range planetListDiv.ChildNodes() {
		// fmt.Printf("child %v - %v\n", child.NodeName(), child.NodeValue())
		planetListDiv.RemoveChild(child)
	}

	planetListDiv.AppendChild(genPlanetHeader(doc))

	planetDisplay = make([]*planetType, 0, len(planetMap))
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

func generatePlanetDisplay(p planetType) dom.Element {

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
	market.SetInnerHTML("Products (my/total amount)")

	category := doc.CreateElement("div")
	category.Class().SetString("planetCategoryMarketHeader")
	category.SetInnerHTML("Category (my/opp share %)")

	headerdiv.AppendChild(info)
	headerdiv.AppendChild(market)
	headerdiv.AppendChild(category)

	return headerdiv
}

func genPlanetInfo(doc dom.Document, p planetType) dom.Element {

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
	editimg.Class().Add("hidden")
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

func genPlanetMarket(doc dom.Document, p planetType) dom.Element {

	marketDiv := doc.CreateElement("div")
	marketDiv.Class().SetString("planetMarket")

	// want to order this correctly.. similar to addPlanetForm
	// this won't be exact, since not all items will be in the list
	// but at least it will try to keep it consistent

	for _, prodId := range productList {
		if prod, exist := p.ProductList[prodId]; exist {
			prodDiv := doc.CreateElement("div")
			prodDiv.SetID(prod.ID())
			prodDiv.SetInnerHTML(fmt.Sprintf("%v:  %v / %v", prod.name, prod.Supply, prod.Demand))
			marketDiv.AppendChild(prodDiv)
		}
	}

	return marketDiv

}

func genPlanetCatMarket(doc dom.Document, p planetType) dom.Element {
	categoryDiv := doc.CreateElement("div")
	categoryDiv.Class().SetString("planetCategoryMarket")
	// var (
	// 	max   int
	// 	maxEl dom.Element
	// )

	// for cat, catMarket := range p.marketByCat {
	for _, cat := range catList {
		if _, exists := p.marketByCat[cat]; exists {

			catMarketDiv := doc.CreateElement("div")
			catMarketDiv.SetID(cat.String())

			if c, o, err := p.calcCategoryShare(cat); err == nil {
				catMarketDiv.SetInnerHTML(fmt.Sprintf("%v: %.1f / %.1f%%", cat.String(), c, o))
			}

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
		var oldDiv = doc.GetElementByID(prevSel).QuerySelector(".planetInfo")
		if oldDiv != nil {
			oldDiv.Class().Remove("selected")
			oldDiv.QuerySelector(".editImage").Class().Add("hidden")
		}
	}

	if newSel == prevSel {
		// same object, its just a deselect.. make sure supply form is hidden
		selected = ""

		supplyForm := doc.GetElementByID("supplyForm")
		if supplyForm.Class().Contains("displayNone") == false {
			supplyForm.Class().Add("displayNone")
		}

		return nil
	}

	// at this point, we're at a new selection.. load supply info
	var newDiv = doc.GetElementByID(newSel).QuerySelector(".planetInfo")
	newDiv.Class().Add("selected")
	newDiv.QuerySelector(".editImage").Class().Remove("hidden")
	selected = newSel

	addPlanetForm := doc.GetElementByID("addPlanetForm")
	if addPlanetForm.Class().Contains("displayNone") == false {
		addPlanetForm.Class().Add("displayNone")
	}

	if err := genSupplyForm(doc, newSel); err != nil {
		return err
	}

	supplyForm := doc.GetElementByID("supplyForm")
	supplyForm.Class().Remove("displayNone")

	return nil
}

// func (p *planet) calcMarketShare() float32 {
// 	return 0
// }

// func getSectorList() []string {
// 	var list = []string{"sector 1", "sector 2", "foobar"}
// 	return list
// }
