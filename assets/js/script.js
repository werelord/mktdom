
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
    if (el.classList.contains("selected")) {
        el.classList.remove('selected')
    } else {
        Array.prototype.slice.call(document.querySelectorAll('div[data-tag="planetList"] div')).forEach(function (element) {
            // remove the selected clas
            element.classList.remove('selected');
        });
        // add the selected class to the element that was clicked
        el.classList.add('selected');
    }

    // todo: update subsequent divs

}

function handleErr(err) {
    console.log("error encocuntered:", err)
    errtxt = document.getElementById("errorText")
    errtxt.innerHTML = err
}

function clearErrTxt() {
    errtxt = document.getElementById("errorText")
    errtxt.innerHTML = ""
}

function addPlanet() {
    console.log("in add planet")
    var result = goAddPlanet()
    if ((result != null) && ('error' in result)) {
        handleErr(result.error)
    }
}


function showAddPlanet() {
    f = document.getElementById("addPlanetForm")
    if ( f.classList.contains("hidden")) {
        f.classList.remove("hidden")
    }
}

function cancelAddPlanet () {
    f = document.getElementById("addPlanetForm")

    document.getElementById("addPlanetName").value = ""
    document.getElementById("addPlanetSector").value = ""

    var result = goGenPlanetForm()
       if ((result != null) && ('error' in result)) {
        handleErr(result.error)
    }
    f.classList.add("hidden")
    clearErrTxt()
}