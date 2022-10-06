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

    if (image in imageUnits) {
        updateRow(image, units, imageUnits[image]);
    } else {
        addRow(image, units);
    }
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

    imageUnits[image] = units;
    console.log(imageUnits);
}

function updateRow(image, units, unitsBefore) {
    const before = document.getElementById("before")
    const after = document.getElementById("after");
    const cancelButton = document.getElementById("cancel");
    const okButton = document.getElementById("ok");

    before.innerText = unitsBefore;
    after.innerText = units;

    dialog.showModal();

    okButton.onclick = () => {
        const unitsTd = document.getElementById(image);
        unitsTd.innerText = units;
        dialog.close();

        imageUnits[image] = units;
        console.log(imageUnits);
    }

    cancelButton.onclick = () => {
        dialog.close();
        console.log(imageUnits);
    }
}

function deleteRow() {
    const removeTr = this.parentNode
    tbody.removeChild(removeTr);

    const image = removeTr.firstChild.innerText;
    delete imageUnits[image];
    console.log(imageUnits);
}

function save() {

    const data = {
        lastName: document.getElementById("lastNameInput").value,
        firsName: document.getElementById("firsNameInput").value,
        nickname: document.getElementById("inputNickname").value,
        mail: document.getElementById("inputMail").value,
        company: document.getElementById("inputCompany").value,
        password: document.getElementById("inputPassword").value,
        userType: document.getElementById("selectUserType").value,
        imageUnits: imageUnits
    }

    register(data);
}

async function register(data) {
    url = 'http://localhost:8080/register'
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    console.log(data);
}