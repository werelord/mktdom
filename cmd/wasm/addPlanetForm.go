package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"honnef.co/go/js/dom/v2"
)

func genAddPlanetForm(planetStr string) any {
	// fmt.Println("todo: " + planetStr)
	var pl planetType

	if planetStr != "" {
		// load planet for editing
		if p, exists := planetMap[planetStr]; exists == false {
			return sendErr("error showing form; planet '%v' does not exist", planetStr)
		} else {
			pl = *p
		}
		// make sure the correct buttons are set
	}
	return genPlanetForm(pl)
}

func genPlanetForm(planet planetType) any {
	// fmt.Println("in genPlanetForm():", p)

	doc := dom.GetWindow().Document()

	// make sure the correct buttons are shown.. keying off the planet name being empty whether its a
	// save or an add
	addButton := doc.GetElementByID("addPlanetButton")
	saveButton := doc.GetElementByID("savePlanetButton")
	deleteButton := doc.GetElementByID("deletePlanetButton")
	if planet.Name == "" {
		addButton.Class().Remove("displayNone")
		saveButton.Class().Add("displayNone")
		deleteButton.Class().Add("displayNone")
	} else {
		addButton.Class().Add("displayNone")
		saveButton.Class().Remove("displayNone")
		deleteButton.Class().Remove("displayNone")

		// save the planet name for delete/save
		saveButton.SetAttribute("data-planet", planet.Name)
		// saveButton.SetAttribute("onclick", fmt.Sprintf("onSavePlanet('%v');", p.Name))
	}

	doc.GetElementByID("addPlanetName").(*dom.HTMLInputElement).SetValue(planet.Name)
	doc.GetElementByID("addPlanetSector").(*dom.HTMLInputElement).SetValue(planet.Sector)
	doc.GetElementByID("addPlanetPoints").(*dom.HTMLInputElement).SetValue(fmt.Sprintf("%v", planet.DomPoints))

	table := doc.GetElementByID("addPlanetTable")
	elemList := table.QuerySelectorAll(".productAmount")

	// fmt.Println("here")
	for _, elem := range elemList {
		var (
			prodName = strings.TrimSuffix(elem.ID(), "_amt")
			prodAmt  int
		)
		if prod, exists := planet.ProductList[prodName]; exists {
			prodAmt = prod.Demand
		} // otherwise zero

		// fmt.Printf("%v:%v\n", prodName, prodAmt)
		elem.(*dom.HTMLInputElement).SetValue(fmt.Sprintf("%v", prodAmt))

	}

	return nil
}

func onAddPlanet(overwrite bool) any {

	// Add or Edit depends on overwrite

	doc := dom.GetWindow().Document()
	var newPlanet planetType
	newPlanet.ProductList = make(map[string]productType, 0)

	saveButton := doc.GetElementByID("savePlanetButton")
	var oldName = saveButton.(*dom.HTMLButtonElement).Dataset()["planet"]
	// fmt.Println("planet : " + oldName)
	if overwrite && (oldName == "") {
		return sendErr("Cannot save planet; cannot find planet with name %v", oldName)
	}

	if name := doc.GetElementByID("addPlanetName").(*dom.HTMLInputElement).Value(); name == "" {
		return sendErr("Add Planet: name cannot be empty")
	} else {

		if overwrite {
			// need to check the old name
			if _, exists := planetMap[oldName]; !exists {
				return sendErr("Unable to find planet entry for old name %v", oldName)
			} else {
				fmt.Printf("changing planet name from %v to %v\n", oldName, name)
			}
		} else {
			if _, exists := planetMap[name]; exists {
				return sendErr("Add Planet: planet with name '%v' already exists", name)
			}
		}
		newPlanet.Name = strings.Title(name)
	}

	if sector := doc.GetElementByID("addPlanetSector").(*dom.HTMLInputElement).Value(); sector == "" {
		// not warning about blank sector
	} else {
		// fmt.Println("sector: " + sector)
		newPlanet.Sector = strings.ToUpper(sector)
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
		if amt := doc.GetElementByID(fmt.Sprintf("%v_amt", prdId)).(*dom.HTMLInputElement).Value(); amt == "" {
			// skip
		} else if val, err := strconv.Atoi(amt); err != nil {
			return sendErr("Add planet: " + err.Error())
		} else if val <= 0 {
			// zero, skip
		} else {
			// force value to be multiple of 48, just in case javascript fails
			if (val % 48) != 0 {
				val = int(math.Round(float64(val)/48.0) * 48)
			}

			fmt.Printf("prod: %v amt: %v\n", prdId, val)
			prod.Demand = val
			newPlanet.ProductList[prdId] = prod
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
		deletePlanetData(oldName)
		if selected == oldName {
			selected = ""
		}
	}
	if err := savePlanetData(newPlanet); err != nil {
		return sendErr("error saving planet data: %v", err)
	}

	/*
		if overwrite remove old, insert new (ignore found)
		if not overwrite, just insert
	*/
	if overwrite {
		if err := removeFromDisplay(doc, &planetType{Name: oldName}); err != nil {
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

		return sendToast("Planet \"%v\" %v", newPlanet.Name, status)
	}
}

func onDeletePlanet() any {

	// fmt.Println("onDeletePlanet()")
	// confirmation should already have happened

	// look in the form for the hidden name that was used; can't use text box (may have been edited)
	// or current selection (may have changed)

	var (
		doc     = dom.GetWindow().Document()
		oldName = doc.GetElementByID("savePlanetButton").(*dom.HTMLButtonElement).Dataset()["planet"]
	)

	if oldName == "" {
		return sendErr("planet name is empty (something is wrong); cannot delete planet")
	} else if planet, exists := planetMap[oldName]; !exists {
		return sendErr("unable to find planet with name %v", oldName)
	} else {
		// fmt.Println("delete planet, name = " + oldName)
		delete(planetMap, oldName)
		deletePlanetData(oldName)
		removeFromDisplay(doc, planet)
		if selected == oldName {
			selected = ""
		}

		return sendToast("Successfully deleted planet '%v'", oldName)
	}
}

func insertIntoDisplay(doc dom.Document, newPlanet planetType) any {
	// do an insert, assuming sorted by name for now
	i, found := slices.BinarySearchFunc(planetDisplay, &newPlanet, func(a, b *planetType) int {
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

func removeFromDisplay(doc dom.Document, remPlanet *planetType) any {
	if i, found := slices.BinarySearchFunc(planetDisplay, remPlanet, func(a, b *planetType) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	}); !found {
		return sendErr("unable to remove planet div '%v'; cannot find div!", remPlanet.Name)
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
