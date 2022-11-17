const selectImage = document.getElementById("image");
const selectUnits = document.getElementById("units");
const tbody = document.getElementById("tbody");

window.addEventListener('DOMContentLoaded', function() {

    initializeImageUnits()

    const editButton = document.getElementById("edit-button")
    const assignButton = document.getElementById("assign-button");
    const updateButton = document.getElementById("update-button");
    const userId = document.getElementById("user-id").value
    const userType = document.getElementById("selectUserType").value

    editButton.addEventListener("click", () => {
        const newPass = prompt("新しいパスワードを入力してください");
        const data = {
            id      : userId,
            password: newPass,
            table: userType,
        }

        updatePass(data).then(message => {
            window.alert(message);
        })
    }, false);

    assignButton.addEventListener("click", assign, false);

    updateButton.addEventListener("click", () => {
        update();
    })
});

async function updatePass(data) {
    url = 'http://localhost:8080/password'
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    
    return res.json();
}

let imageUnits = {};
function initializeImageUnits() {
    TrImageUnits = document.getElementsByClassName("image-units");
    for (var tr of TrImageUnits) {
        const tdFileName = tr.getElementsByClassName("file-name");
        const tdLotteryUnits = tr.getElementsByClassName("lottery-units");
        fileName = tdFileName[0].innerText;
        imageId = tdFileName[0].getAttribute('value');
        lotteryUnits = tdLotteryUnits[0].innerText;

        imageUnits[fileName] = [imageId, lotteryUnits]

        // set delete button
        const deleteTd = document.createElement("td");
        deleteTd.classList.add("bi", "bi-trash");
        deleteTd.addEventListener("click", deleteRow, false);
        tr.appendChild(deleteTd);
    }
}

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


function update() {
    const data = {
        id: document.getElementById("user-id").value,
        familyName: document.getElementById("familyNameInput").value,
        firstName: document.getElementById("firstNameInput").value,
        nickname: document.getElementById("inputNickname").value,
        mail: document.getElementById("inputMail").value,
        company: document.getElementById("inputCompany").value,
        userType: document.getElementById("selectUserType").value,
        imageUnits: imageUnits,
    }
    send(data);
}

async function send(data) {
    url = 'http://localhost:8080/edit'
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    const resObj = await res.json();

    if (resObj["invalid"]) {
        invalidMessage(resObj["resmap"]);
    } else {
        window.alert("登録が完了しました");
    }
}

function invalidMessage(resMap) {
    const mailAlert = document.getElementById("mail-alert");
    const nicknameAlert = document.getElementById("nickname-alert")
    if (resMap["mail"] == 0) {
        mailAlert.innerText = "このメールアドレスはすでに登録されています"
    } else {
        mailAlert.innerText = ""
    }
    if (resMap["nickname"] == 0) {
        nicknameAlert.innerText = "このニックネームはすでに登録されています"
    } else {
        nicknameAlert.innerText = ""
    }
}