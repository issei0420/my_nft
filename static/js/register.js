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

    const idx = selectImage.selectedIndex
    const fileName = selectImage.options[idx].text;
    const imageId = selectImage.value;
    const units = selectUnits.value;

    if (fileName === "" || units === "") {
        return
    }

    if (fileName in imageUnits) {
        updateRow(fileName, imageId, units);
    } else {
        addRow(fileName, imageId, units);
    }
}

function addRow(fileName, imageId, units) {
    const imageTd = document.createElement("td");
    imageTd.innerText = fileName;
    const unitsTd = document.createElement("td");
    unitsTd.innerText = units
    unitsTd.setAttribute("id", fileName);
    const deleteTd = document.createElement("td");
    deleteTd.classList.add("bi", "bi-trash");
    deleteTd.addEventListener("click", deleteRow, false);

    const tr = document.createElement("tr");
    tr.appendChild(imageTd);
    tr.appendChild(unitsTd);
    tr.appendChild(deleteTd);

    tbody.appendChild(tr);

    imageUnits[fileName] = [imageId, units];
}

function updateRow(fileName, imageId, units) {
    const before = document.getElementById("before")
    const after = document.getElementById("after");
    const cancelButton = document.getElementById("cancel");
    const okButton = document.getElementById("ok");

    before.innerText = imageUnits[fileName][1];
    after.innerText = units;

    dialog.showModal();

    okButton.onclick = () => {
        const unitsTd = document.getElementById(fileName);
        unitsTd.innerText = units;
        dialog.close();

        imageUnits[fileName] = [imageId, units];
    }

    cancelButton.onclick = () => {
        dialog.close();
    }
}

function deleteRow() {
    const removeTr = this.parentNode
    tbody.removeChild(removeTr);

    const fileName = removeTr.firstChild.innerText;
    delete imageUnits[fileName];
}

function save() {

    res = validate();
    if (res == "invalid") {
        return
    }

    const data = {
        familyName: document.getElementById("familyNameInput").value,
        firsName: document.getElementById("firsNameInput").value,
        nickname: document.getElementById("inputNickname").value,
        mail: document.getElementById("inputMail").value,
        company: document.getElementById("inputCompany").value,
        password: document.getElementById("inputPassword").value,
        userType: document.getElementById("selectUserType").value,
        imageUnits: imageUnits,
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
    const resObj = await res.json();

    if (resObj["invalid"]) {
        invalidMessage(resObj["resmap"]);
    } else {
        window.alert("???????????????????????????")
    }
}

function invalidMessage(resMap) {
    const mailAlert = document.getElementById("mail-alert");
    const nicknameAlert = document.getElementById("nickname-alert")
    if (resMap["mail"] == 0) {
        mailAlert.innerText = "???????????????????????????????????????????????????????????????"
    } else {
        mailAlert.innerText = ""
    }
    if (resMap["nickname"] == 0) {
        nicknameAlert.innerText = "????????????????????????????????????????????????????????????"
    } else {
        nicknameAlert.innerText = ""
    }
}


// ?????????????????????????????????(validate.js????????????)

const theForm = document.getElementById("form")
const requiredElems = document.querySelectorAll('.required');
const errorClassName = 'error';

// ???????????????????????????????????????????????????????????????
const createError = (elem, errorMessage) => {
    const errorSpan = document.createElement('span');
    errorSpan.classList.add(errorClassName);
    errorSpan.setAttribute('aria-live', 'polite');
    errorSpan.textContent = errorMessage;
    errorSpan.style.color = "#dc3545";
    elem.parentNode.appendChild(errorSpan);
}    

function validate() {
    //?????????
    const errorElems = theForm.querySelectorAll('.' + errorClassName);
    errorElems.forEach( (elem) => {
        elem.remove(); 
    });

    //.required ??????????????????????????????
    let res = ""
    requiredElems.forEach((elem) => {
        const elemValue = elem.value.trim(); 
        if(elemValue.length === 0) {
            createError(elem, '?????????????????????');
            res = "invalid";
        }
    });

    return res
}

// ??????????????????????????????????????????????????????????????????
const selectUserType = document.getElementById("selectUserType")
selectUserType.addEventListener("change", (e) => {
    if (selectUserType.value == "sellers") {
      selectImage.disabled = true;
      selectImage.value = "";
      selectUnits.disabled = true;
      selectUnits.value = ""
    } else {
      selectImage.disabled = false;
      selectUnits.disabled = false;
    }
    e.preventDefault();
});