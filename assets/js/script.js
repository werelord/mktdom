
function showToast (msg) {
    var toast = document.getElementById("mtoast")
    toast.innerHTML = msg
    toast.classList.add("show")

    setTimeout(function() { toast.className = toast.classList.remove("show"); }, 2500)
}

function addAmount(el) {
    var e = document.getElementById(el)
    e.innerHTML = +e.innerHTML + 48
}

function subAmount(el) {
    var e = document.getElementById(el)
    if (+e.innerHTML > 0) {
        e.innerHTML = +e.innerHTML - 48
    }
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
        } else if ('toast' in result) {
            // make sure any errors are cleared
            errtxt.classList.remove("show")
            showToast(result.toast)
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


function showAddPlanet() {
    handleResult(goShowAddPlanet(""))
    f = document.getElementById("addPlanetForm")
    if ( f.classList.contains("hidden")) {
        f.classList.remove("hidden")
    }
}

function cancelAddPlanet () {
    f = document.getElementById("addPlanetForm")

    clearAddPlanet()
    clearErrTxt()
    f.classList.add("hidden")
}

function clearAddPlanet() {
    f = document.getElementById("addPlanetForm")

    document.getElementById("addPlanetName").value = ""
    document.getElementById("addPlanetSector").value = ""
    document.getElementById("addPlanetPoints").value = ""
}

function editPlanet(planet, event) {
    // console.log("edit planet " + planet)
    // console.log("event " + event)
    handleResult(goShowAddPlanet(planet))
    event.stopPropagation()

}