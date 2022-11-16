 
const selectUnits = document.getElementById("selectUnits");

const selectImage = document.getElementById("selectImage");

selectImage.addEventListener("change", () => {
    selected = selectImage.options[selectImage.selectedIndex]
    units = Number(selected.getAttribute("class"))
    expandSelection(units);
})

function expandSelection(units) {
    initializeSelectOption()
    for (let i = 1; i < units + 1; i++) {
        optionUnits = document.createElement("option");
        optionUnits.innerHTML = i;
        selectUnits.appendChild(optionUnits);
    }
}

function initializeSelectOption() {
    while(selectUnits.firstChild){
        selectUnits.removeChild(selectUnits.firstChild);
    }
    selected = document.createElement("option");
    selected.innerHTML = "未選択"
    selected.setAttribute('selected', '');
    selectUnits.appendChild(selected);
}