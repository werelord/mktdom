
function showInfoToast(msg) {
    var toast = document.getElementById("infoToast")
    toast.innerHTML = msg
    toast.classList.add("show")

    setTimeout(function () { toast.classList.remove("show"); }, 2500)
}

function showErrToast(msg) {
    var toast = document.getElementById("errToast")
    toast.innerHTML = msg
    toast.classList.add("show")

    setTimeout(function () { toast.classList.remove("show"); }, 2500)
}

function addAmount(el) {
    var e = document.getElementById(el)
    e.innerHTML = +e.innerHTML + 48
}

function subAmount(el) {
    var e = document.getElementById(el)
    var amt = +e.innerHTML - 48
    if (amt < 0) {
        amt = 0
    }
    e.innerHTML = amt
}

function switchSelected(el) {
    handleResult(goOnSelected(el))
    // if (el.classList.contains("selected")) {
    //     el.classList.remove('selected')
    // } else {
    //     Array.prototype.slice.call(document.querySelectorAll('div[data-tag="planetList"] div')).forEach(function (element) {
    //         // remove the selected clas
    //         element.classList.remove('selected');
    //     });
    //     // add the selected class to the element that was clicked
    //     el.classList.add('selected');
    // }

    // todo: update subsequent divs

}

function handleResult(result) {
    if (result != null) {
        errtxt = document.getElementById("errorText")
        if ('error' in result) {
            console.log("error encocuntered:", result.error)
            errtxt.innerHTML = result.error
            errtxt.classList.add("show")
            showErrToast(result.error)
        } else if ('toast' in result) {
            // make sure any errors are cleared
            errtxt.classList.remove("show")
            showInfoToast(result.toast)
            return true
        }
    }
}

function clearErrTxt() {
    errtxt = document.getElementById("errorText")
    // errtxt.innerHTML = ""
    errtxt.classList.remove("show")
}

function onAddPlanet() {
    console.log("in add planet")
    if (handleResult(goOnAddPlanet())) {
        clearAddPlanet()
    }
}

function onSavePlanet() {
    console.log("in save planet")
    if (handleResult(goOnSavePlanet())) {
        clearAddPlanet()
        hideAddPlanet()
    }
}

function onDeletePlanet() {
    // confirmation handled by dialog
    console.log("in delete planet")
    if (handleResult(goOnDeletePlanet())) {
        clearAddPlanet()
    }
}

function showAddPlanet() {
    handleResult(goGenAddPlanetForm(""))
    supplyForm = document.getElementById("supplyForm")
    if (supplyForm.classList.contains("displayNone") == false) {
        supplyForm.classList.add("displayNone")
    }
    addForm = document.getElementById("addPlanetForm")
    if (addForm.classList.contains("displayNone")) {
        addForm.classList.remove("displayNone")
    }
}

function hideAddPlanet() {
    addForm = document.getElementById("addPlanetForm")
    if (addForm.classList.contains("displayNone") == false) {
        addForm.classList.add("displayNone")
    }
    clearAddPlanet()
    clearErrTxt()
}

function clearAddPlanet() {
    // no reason to call into go for this
    f = document.getElementById("addPlanetForm")

    document.getElementById("addPlanetName").value = ""
    document.getElementById("addPlanetSector").value = ""
    document.getElementById("addPlanetPoints").value = ""
}

function editPlanet(planet, event) {
    handleResult(goGenAddPlanetForm(planet))
    hideSupplyForm()
    f = document.getElementById("addPlanetForm")
    if (f.classList.contains("displayNone")) {
        f.classList.remove("displayNone")
    }
    event.stopPropagation()
}

function hideSupplyForm() {
    f = document.getElementById("supplyForm")
    if (f.classList.contains("displayNone") == false) {
        f.classList.add("displayNone")
    }
}

function changeSupply(op, planet, product) {
    handleResult(goOnChangeSupply(op, planet, product))
}
