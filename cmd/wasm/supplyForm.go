package main

import (
	"fmt"
	"math"

	"honnef.co/go/js/dom/v2"
)

func genSupplyForm(doc dom.Document, planetStr string) any {

	var (
		p        planetType
		title    = doc.GetElementByID("supplyTitle")
		myShare  = doc.GetElementByID("infoMyShare")
		oppShare = doc.GetElementByID("infoOppShare")
		// prodTable = doc.GetElementByID("supplyFormProductTable")
	)

	if planet, exists := planetMap[planetStr]; !exists {
		return sendErr("unable to find planet with name '%v'", planetStr)
	} else {
		p = *planet
	}

	title.SetInnerHTML(fmt.Sprintf("Planet \"%v\": Edit Supply", p.Name))
	title.SetAttribute("data-planet", p.Name)
	myShare.SetInnerHTML(fmt.Sprintf("My Share (total): %.2f%%",
		float32(p.market.current)/float32(p.market.total)*100.0))

	cat, share := p.calcMaxOppShare()
	oppShare.SetInnerHTML(fmt.Sprintf("Max Opponent Share (%v): %.2f%%", cat.String(), share))

	// fmt.Printf("genSupplyForm, %v\n", planetStr)

	genSupplyCategoryTable(doc, p)

	genSupplyProductTable(doc, p)

	// for _, child := range prodTable.ChildNodes() {
	// 	prodTable.RemoveChild(child)
	// }

	return nil
	// return sendErr("not yet implemented")
}

func genSupplyCategoryTable(doc dom.Document, p planetType) {

	// fmt.Printf("genSupplyCategoryTalble: %v\n", p.Name)

	var catTable = doc.GetElementByID("supplyFormCategoryTable")
	for _, child := range catTable.ChildNodes() {
		// fmt.Println("removing node")
		catTable.RemoveChild(child)
	}
	// fmt.Println("wtf2")

	thead := doc.CreateElement("thead")
	catTable.AppendChild(thead)

	mine := doc.CreateElement("td")
	mine.SetInnerHTML("Mine")
	opp := doc.CreateElement("td")
	opp.SetInnerHTML("Opp")
	hid := doc.CreateElement("td")
	hid.Class().SetString("hidden")
	hid.SetInnerHTML("0000")

	thead.AppendChild(doc.CreateElement("td"))
	thead.AppendChild(mine)
	thead.AppendChild(opp)
	if len(p.marketByCat) > 1 {
		thead.AppendChild(hid)
		thead.AppendChild(doc.CreateElement("td"))
		thead.AppendChild(mine.CloneNode(true))
		thead.AppendChild(opp.CloneNode(true))
	}

	tbody := doc.CreateElement("tbody")
	catTable.AppendChild(tbody)

	var (
		tr    dom.Element
		start = true
	)
	for _, cat := range catList {

		if _, exists := p.marketByCat[cat]; exists {
			// iterate left of table, right of table
			if start {
				tr = doc.CreateElement("tr")
				tbody.AppendChild(tr)
				start = false
			} else {
				tr.AppendChild(hid.CloneNode(true))
				start = true
			}

			name := doc.CreateElement("td")
			mine := doc.CreateElement("td")
			opp := doc.CreateElement("td")

			name.SetID(fmt.Sprintf("sf_%v", cat.String()))
			name.SetInnerHTML(cat.String())
			var m, o, _ = p.calcCategoryShare(cat)
			mine.SetInnerHTML(fmt.Sprintf("%.2f%%", m))
			opp.SetInnerHTML(fmt.Sprintf("%.2f%%", o))

			tr.AppendChild(name)
			tr.AppendChild(mine)
			tr.AppendChild(opp)
		}
	}

	// return nil
}

func genSupplyProductTable(doc dom.Document, p planetType) {
	fmt.Printf("genSupplyProdTable: %v\n", p.Name)

	var prodTable = doc.GetElementByID("supplyFormProductTable")

	for _, child := range prodTable.ChildNodes() {
		prodTable.RemoveChild(child)
	}

	tbody := doc.CreateElement("tbody")
	prodTable.AppendChild(tbody)

	// hidden
	hidden := doc.CreateElement("tr")
	hidden.Class().SetString("hidden")
	tbody.AppendChild(hidden)

	spacer := doc.CreateElement("td")
	spacer.SetInnerHTML("0000")

	hidden.AppendChild(doc.CreateElement("td"))
	hidden.AppendChild(doc.CreateElement("td"))
	hidden.AppendChild(spacer)
	hidden.AppendChild(doc.CreateElement("td"))
	if len(p.ProductList) > 1 {
		hidden.AppendChild(doc.CreateElement("td"))
		hidden.AppendChild(doc.CreateElement("td"))
		hidden.AppendChild(spacer.CloneNode(true))
		hidden.AppendChild(doc.CreateElement("td"))
	}

	var (
		tr    dom.Element
		start = true
	)

	for _, prodStr := range productList {
		if prod, exists := p.ProductList[prodStr]; exists {
			if start {
				tr = doc.CreateElement("tr")
				tbody.AppendChild(tr)
				start = false
			} else {
				start = true
			}
			var perChange = (48.0 * float64(prod.price) / float64(p.market.total) * 100)
			perChange = math.Round(perChange*100) / 100

			name := doc.CreateElement("td")
			name.SetInnerHTML(fmt.Sprintf("%v (&#177;%.1f%%)", prod.name, perChange))

			sub := doc.CreateElement("td")
			subButton := doc.CreateElement("button")
			subButton.SetAttribute("onclick", fmt.Sprintf("subSupply('%v')", prod.ID()))
			subButton.SetInnerHTML("-")
			sub.AppendChild(subButton)

			amt := doc.CreateElement("td")
			amt.SetID(fmt.Sprintf("%v_supply", prod.ID()))
			amt.SetInnerHTML(fmt.Sprintf("%v", prod.Supply))

			add := doc.CreateElement("td")
			addButton := doc.CreateElement("button")
			addButton.SetAttribute("onclick", fmt.Sprintf("addSupply('%v')", prod.ID()))
			addButton.SetInnerHTML("+")
			add.AppendChild(addButton)

			tr.AppendChild(name)
			tr.AppendChild(sub)
			tr.AppendChild(amt)
			tr.AppendChild(add)
		}
	}
}

func onChangeSupply(amt int, productStr string) any {

	return sendErr("not yet implemented")
}
