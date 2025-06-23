package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"honnef.co/go/js/dom/v2"
)

func showAddPlanet(planetStr string) any {
	// fmt.Println("todo: " + planetStr)
	var pl planet

	if planetStr != "" {
		// load planet for editing
		if p, exists := planetMap[planetStr]; exists == false {
			return sendErr(fmt.Sprintf("error showing form; planet '%v' does not exist", planetStr))
		} else {
			pl = *p
		}
		// make sure the correct buttons are set
	}
	return genPlanetForm(pl)
}

func genPlanetForm(p planet) any {
	// fmt.Println("in genPlanetForm():", p)

	doc := dom.GetWindow().Document()

	// make sure the correct buttons are shown.. keying off the planet name being empty whether its a
	// save or an add
	addButton := doc.GetElementByID("addPlanetButton")
	saveButton := doc.GetElementByID("savePlanetButton")
	if p.Name == "" {
		addButton.Class().Remove("displayNone")
		saveButton.Class().Add("displayNone")
	} else {
		addButton.Class().Add("displayNone")
		saveButton.Class().Remove("displayNone")
	}

	doc.GetElementByID("addPlanetName").(*dom.HTMLInputElement).SetValue(p.Name)
	doc.GetElementByID("addPlanetSector").(*dom.HTMLInputElement).SetValue(p.Sector)
	doc.GetElementByID("addPlanetPoints").(*dom.HTMLInputElement).SetValue(fmt.Sprintf("%v", p.DomPoints))

	table := doc.GetElementByID("addPlanetTable")

	// reset everything
	for _, child := range table.ChildNodes() {
		table.RemoveChild(child)
	}

	// add spacer row
	hiddenrow := doc.CreateElement("tr")
	hiddenrow.Class().SetString("hidden")
	// save original name, just in case name changes
	hiddenName := doc.CreateElement("td")
	spacer1 := doc.CreateElement("td")
	spacer2 := doc.CreateElement("td")
	hiddenName.SetID("old_name")
	hiddenName.SetInnerHTML(p.Name)
	spacer1.SetInnerHTML("0000")
	spacer2.SetInnerHTML("0000")

	hiddenrow.AppendChild(hiddenName)
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(spacer1)
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(doc.CreateElement("td"))
	hiddenrow.AppendChild(spacer2)
	hiddenrow.AppendChild(doc.CreateElement("td"))

	table.AppendChild(hiddenrow)

	var genProdTd = func(root dom.Element, p planet, prod productType) {
		// var tdlist = make([]dom.Element, 0)

		label := doc.CreateElement("td")
		label.Class().SetString("productLabel")
		// label.SetAttribute("id", fmt.Sprintf("%v", spToUl(p.ID())))
		// fmt.Println("in genPlanetForm, " + p.name)
		label.SetInnerHTML(prod.name)
		root.AppendChild(label)

		mbntd := doc.CreateElement("td")
		minus := doc.CreateElement("button")
		minus.Class().SetString("productButton")
		minus.SetAttribute("onclick", fmt.Sprintf("subAmount('%v_amt')", prod.ID()))
		minus.SetInnerHTML("-")
		mbntd.AppendChild(minus)
		root.AppendChild(mbntd)

		amt := doc.CreateElement("td")
		amt.Class().SetString("productAmount")
		amt.SetID(fmt.Sprintf("%v_amt", prod.ID()))
		var prodAmt int

		// fmt.Printf("planet %v, product %v\n", p, prod.ID())
		if existingProd, exists := p.ProductList[prod.ID()]; exists {
			// fmt.Printf("exists!! %v\n", existingProd)
			prodAmt = existingProd.Demand
		}
		amt.SetInnerHTML(fmt.Sprintf("%v", prodAmt))
		root.AppendChild(amt)

		pbntd := doc.CreateElement("td")
		plus := doc.CreateElement("button")
		plus.Class().SetString("productButton")
		plus.SetAttribute("onclick", fmt.Sprintf("addAmount('%v_amt')", prod.ID()))
		plus.SetInnerHTML("+")
		pbntd.AppendChild(plus)
		root.AppendChild(pbntd)

	}

	for i := 0; i < len(productList); i = i + 2 {
		var p1, p2 = baseProductMap[productList[i]], baseProductMap[productList[i+1]]

		row := doc.CreateElement("tr")

		genProdTd(row, p, p1)
		genProdTd(row, p, p2)

		// plus := doc.CreateElement("button")

		table.AppendChild(row)
	}

	return nil
}

func onAddPlanet(overwrite bool) any {

	// Add or Edit depends on overwrite

	doc := dom.GetWindow().Document()
	var newPlanet planet
	newPlanet.ProductList = make(map[string]productType, 0)

	var oldName = doc.GetElementByID("old_name").InnerHTML()
	if overwrite && (oldName == "") {
		return sendErr("Cannot save planet; cannot find old planet name")
	}

	if name := doc.GetElementByID("addPlanetName").(*dom.HTMLInputElement).Value(); name == "" {
		return sendErr("Add Planet: name cannot be empty")
	} else {

		if overwrite {
			// need to check the old name
			if _, exists := planetMap[oldName]; !exists {
				return sendErr(fmt.Sprintf("Unable to find planet entry for old name %v", oldName))
			} else {
				fmt.Printf("changing planet name from %v to %v\n", oldName, name)
			}
		} else {
			if _, exists := planetMap[name]; exists {
				return sendErr(fmt.Sprintf("Add Planet: planet with name '%v' already exists", name))
			}
		}
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
	// make sure old planet is deleted if overwrite and names differ
	if overwrite && (oldName != newPlanet.Name) {
		delete(planetMap, oldName)
		if selected == oldName {
			selected = ""
		}
	}
	if err := savePlanetData(); err != nil {
		return sendErr(fmt.Sprintf("error saving planet data: %v", err))
	}

	/*
		if overwrite remove old, insert new (ignore found)
		if not overwrite, just insert
	*/
	if overwrite {
		if err := removeFromDisplay(doc, &planet{Name: oldName}); err != nil {
			return err
		}
	}

	if err := insertIntoDisplay(doc, newPlanet); err != nil {
		return err
	} else {
		var status = "added"
		if overwrite {
			status = "updated"
		}

		return sendToast(fmt.Sprintf("Planet \"%v\" %v", newPlanet.Name, status))
	}
}

func insertIntoDisplay(doc dom.Document, newPlanet planet) any {
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
	return nil
}

func removeFromDisplay(doc dom.Document, remPlanet *planet) any {
	if i, found := slices.BinarySearchFunc(planetDisplay, remPlanet, func(a, b *planet) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	}); !found {
		return sendErr(fmt.Sprintf("unable to remove planet div '%v'; cannot find div!", remPlanet.Name))
	} else {
		// fmt.Printf("before %v\n", planetDisplay)
		planetDisplay = append(planetDisplay[:i], planetDisplay[i+1:]...)
		// fmt.Printf("after  %v\n", planetDisplay)
	}

	pl := doc.GetElementByID("planetList")
	remDiv := doc.GetElementByID(fmt.Sprintf(remPlanet.Name))
	pl.RemoveChild(remDiv)

	return nil
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
