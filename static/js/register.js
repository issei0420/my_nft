const assignButton = document.getElementById("assign-button");
const saveButton = document.getElementById("save-button")
const selectImage = document.getElementById("image");
const selectUnits = document.getElementById("units");
const tbody = document.getElementById("tbody");

window.addEventListener("DOMContentLoaded", () => {

    assignButton.addEventListener("click", assign, false);
    saveButton.addEventListener("click", save, false);

}, false);

let imageUnits = {};

function assign() {

    const image = selectImage.value;
    const units = selectUnits.value;

    if (image === "" || units === "") {
        return
    }

    // 画面を更新
    if (image in imageUnits) {
        updateRow(image, units);
    } else {
        addRow(image, units);
    }

    // 変数を更新
    imageUnits[image] = units;
    console.log(imageUnits);
}

function addRow(image, units) {
    const imageTd = document.createElement("td");
    imageTd.innerText = image;
    const unitsTd = document.createElement("td");
    unitsTd.innerText = units
    unitsTd.setAttribute("id", image);
    const deleteTd = document.createElement("td");
    deleteTd.classList.add("bi", "bi-trash");
    deleteTd.addEventListener("click", deleteRow, false);

    const tr = document.createElement("tr");
    tr.appendChild(imageTd);
    tr.appendChild(unitsTd);
    tr.appendChild(deleteTd);

    tbody.appendChild(tr);
}

function updateRow(image, units) {
    const unitsTd = document.getElementById(image);
    unitsTd.innerText = units;
}

function deleteRow() {
    const removeTr = this.parentNode
    tbody.removeChild(removeTr);

    const image = removeTr.firstChild.innerText;
    delete imageUnits[image];
    console.log(imageUnits);
}

function save() {
    console.log("確認");
}