package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"honnef.co/go/js/dom/v2"
)

func genPlanetForm() any {
	// fmt.Println("in genPlanetForm()")

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

func addPlanet() any {

	doc := dom.GetWindow().Document()
	var newPlanet planet
	newPlanet.ProductList = make(map[string]productType, 0)

	if name := doc.GetElementByID("addPlanetName").(*dom.HTMLInputElement).Value(); name == "" {
		return sendErr("Add Planet: name cannot be empty")
	} else if _, exists := planetMap[name]; exists {
		return sendErr(fmt.Sprintf("Add Planet: planet with name '%v' already exists", name))
	} else {
		// fmt.Println("name:" + nameEl.Value())
		newPlanet.Name = name
	}

	if sector := doc.GetElementByID("addPlanetSector").(*dom.HTMLInputElement).Value(); sector == "" {
		// not warning about blank sector
	} else {
		// fmt.Println("sector: " + sector)
		newPlanet.Sector = sector
	}

	if domPts := doc.GetElementByID("addPlanetPoints").(*dom.HTMLInputElement).Value(); domPts == "" {
		return sendErr("Add Planet: domination points should not be empty")
	} else if val, err := strconv.Atoi(domPts); err != nil {
		return sendErr("Add Planet: unable to convert domination points to integer")
	} else if val <= 0 {
		return sendErr("Add Planet: domination points should be positive")
	} else {
		newPlanet.DomPoints = val
	}

	for prdId, prod := range baseProductMap {
		if amt := doc.GetElementByID(fmt.Sprintf("%v_amt", prdId)).InnerHTML(); amt == "" {
			// skip
		} else if val, err := strconv.Atoi(amt); err != nil {
			return sendErr("Add planet: " + err.Error())
		} else if val <= 0 {
			// zero, skip
		} else {
			fmt.Printf("prod: %v amt: %v\n", prdId, val)
			prod.Demand = val
			newPlanet.ProductList[prdId] = prod

			// vol := prod.price * prod.Demand
			// planet.currentMarketVol += vol
			// planet.currentMarketVolByCat[prod.category] += vol
		}
	}

	if len(newPlanet.ProductList) == 0 {
		return sendErr("Add Planet: no products defined; must include at least one product")
	}

	newPlanet.calcMarketVol()
	// fmt.Printf("adding planet: %+v\n", newPlanet)

	planetMap[newPlanet.Name] = &newPlanet
	if err := savePlanetData(); err != nil {
		return sendErr(fmt.Sprintf("error saving planet data: %v", err))
	}

	// do an insert, assuming sorted by name for now
	i, found := slices.BinarySearchFunc(planetDisplay, &newPlanet, func(a, b *planet) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))

	})
	if found {
		// this shouldn't happen, name was already checked
		return sendErr("cannot add; name already exists (this shouldn't happen here)")
	}

	newPlanetDiv := generatePlanetDisplay(newPlanet)
	pl := doc.GetElementByID("planetList")

	if i == len(planetDisplay) {
		planetDisplay = append(planetDisplay, &newPlanet)
		// insert at the end
		pl.AppendChild(newPlanetDiv)

	} else {
		// save current planet for dom insertion
		var beforePlanet = planetDisplay[i]

		planetDisplay = append(planetDisplay[:i+1], planetDisplay[i:]...)
		planetDisplay[i] = &newPlanet

		// insert before current i
		befDiv := doc.GetElementByID(beforePlanet.Name)
		pl.InsertBefore(newPlanetDiv, befDiv)
	}

	return sendToast(fmt.Sprintf("Planet \"%v\" added", newPlanet.Name))
}

func (p *planet) calcMarketVol() {
	// fmt.Printf("before planet: %+v\n", *p)
	// zero out before calc, just in case this is called elsewhere
	p.market.current = 0
	p.market.total = 0
	// p.market.share = 0.0
	p.marketByCat = make(map[categoryType]marketVolume, 0)

	var (
		totVol int
		curVol int
	)

	for _, prod := range p.ProductList {
		curVol = (prod.Supply * prod.price)
		totVol = (prod.Demand * prod.price)

		p.market.current += curVol
		p.market.total += totVol

		var (
			mkt    marketVolume
			exists bool
		)

		if mkt, exists = p.marketByCat[prod.category]; exists == false {
			mkt = marketVolume{}
			// p.marketByCat[prod.category] = mkt
		}
		mkt.current += curVol
		mkt.total += totVol
		p.marketByCat[prod.category] = mkt

	}
	// fmt.Printf("after planet: %+v\n", *p)
}
