package main

import (
	"fmt"
	"sort"
	"strings"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

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
	js.Global().Set("goGenPlanetForm", js.FuncOf(func(this js.Value, args []js.Value) any {
		return genPlanetForm()
	}))

	<-make(chan struct{})
}

func loadStoredPlanetData() {
	// todo: load shit
	temploadStoredPlanetData()
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
		fullPlanetDiv.Class().SetString("fullPlanetDetails")
		fullPlanetDiv.SetAttribute("onClick", "switchSelected(this)")

		planetWrapperDiv := doc.CreateElement("div")
		// planetWrapperDiv.SetAttribute("onClick", "switchSelected(this)")
		fullPlanetDiv.AppendChild(planetWrapperDiv)

		nameDiv := doc.CreateElement("div")
		nameDiv.Class().SetString("planetName")
		nameDiv.SetInnerHTML(planet.name)
		planetWrapperDiv.AppendChild(nameDiv)

		detailsDiv := doc.CreateElement("div")
		detailsDiv.Class().SetString("planetDetails")
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
		marketDiv.Class().SetString("planetMarket")
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
	fmt.Println("in genPlanetForm()")
	// return nil
	doc := dom.GetWindow().Document()
	table := doc.GetElementByID("addPlanetTable")

	// reset everything
	for _, child := range table.ChildNodes() {
		// fmt.Printf("child %v - %v\n", child.NodeName(), child.NodeValue())
		table.RemoveChild(child)
	}

	// add spacer row
	hiddenrow := doc.CreateElement("tr")
	hiddenrow.Class().SetString("hidden")
	spacer1 := doc.CreateElement("td")
	spacer2 := doc.CreateElement("td")
	spacer1.SetInnerHTML("0000")
	spacer2.SetInnerHTML("0000")

	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(spacer1)
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(spacer2)
	hiddenrow.AppendChild(doc.CreateElement("td"))

	table.AppendChild(hiddenrow)

	var genProdTd = func(r dom.Element, p productType) {
		// var tdlist = make([]dom.Element, 0)

		label := doc.CreateElement("td")
		label.Class().SetString("productLabel")
		// label.SetAttribute("id", fmt.Sprintf("%v", spToUl(p.ID())))
		// fmt.Println("in genPlanetForm, " + p.name)
		label.SetInnerHTML(p.name)
		r.AppendChild(label)

		mbntd := doc.CreateElement("td")
		minus := doc.CreateElement("button")
		minus.Class().SetString("productButton")
		minus.SetAttribute("onclick", fmt.Sprintf("subAmount('%v_amt')", p.ID()))
		minus.SetInnerHTML("-")
		mbntd.AppendChild(minus)
		r.AppendChild(mbntd)

		amt := doc.CreateElement("td")
		amt.Class().SetString("productAmount")
		amt.SetID(fmt.Sprintf("%v_amt", p.ID()))
		amt.SetInnerHTML(fmt.Sprintf("%v", 0))
		r.AppendChild(amt)

		pbntd := doc.CreateElement("td")
		plus := doc.CreateElement("button")
		plus.Class().SetString("productButton")
		plus.SetAttribute("onclick", fmt.Sprintf("addAmount('%v_amt')", p.ID()))
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
