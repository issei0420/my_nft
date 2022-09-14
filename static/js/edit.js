const editButton = document.getElementById("edit-button")
editButton.addEventListener("click", () => {
    let newPass = prompt("新しいパスワードを入力してください");
    updatePass(newPass).then(message => {
        console.log(message);
    })
}, false);

async function updatePass(newPass) {
    console.log(newPass);
    url = 'http://localhost:8080/password'
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newPass)
    })
    
    return res.json();
}
