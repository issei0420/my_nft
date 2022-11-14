window.addEventListener('DOMContentLoaded', function() {

    initializeImageUnits()

    const editButton = document.getElementById("edit-button")
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
async function initializeImageUnits() {
    TrImageUnits = document.getElementsByClassName("image-units");
    for (var tr of TrImageUnits) {
        const tdFileName = tr.getElementsByClassName("file-name");
        const tdLotteryUnits = tr.getElementsByClassName("lottery-units");
        fileName = tdFileName[0].innerText;
        imageId = tdFileName[0].getAttribute('value');
        lotteryUnits = tdLotteryUnits[0].innerText;

        imageUnits[fileName] = [imageId, lotteryUnits]
    }
    console.log(imageUnits);
}
