const assignButton = document.getElementById("assign-button");
const saveButton = document.getElementById("save-button")
const selectImage = document.getElementById("image");
const selectUnits = document.getElementById("units");
const tbody = document.getElementById("tbody");

window.addEventListener("DOMContentLoaded", () => {

    assignButton.addEventListener("click", assign, false);
    saveButton.addEventListener("click", save, false);

}, false);


function assign() {

    const image = selectImage.value;
    const units = selectUnits.value;

    if (image === "" || units === "") {
        return
    }

    const imageTd = document.createElement("td");
    imageTd.innerText = image;
    const UnitsTd = document.createElement("td");
    UnitsTd.innerText = units
    const deleteTd = document.createElement("td");
    deleteTd.classList.add("bi", "bi-trash");
    deleteTd.addEventListener("click", deleteRow, false);

    const tr = document.createElement("tr");
    tr.appendChild(imageTd);
    tr.appendChild(UnitsTd);
    tr.appendChild(deleteTd);

    tbody.appendChild(tr);
}

function deleteRow() {
    removeTr = this.parentNode
    tbody.removeChild(removeTr);
}

function save() {
    console.log("確認");
}