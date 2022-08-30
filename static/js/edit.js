window.addEventListener("DOMContentLoaded", () => {
    const yaunValue = document.getElementById("selectUserType").value 
    initialize(yaunValue);
}, false)

function initialize(yaunValue) {
    const yaunTable = document.createElement("input");
    yaunTable.type = "hidden";
    yaunTable.name = "yaunTable"
    yaunTable.id = "yaunTable"
    yaunTable.value = yaunValue

    const form = document.getElementById("form");
    form.appendChild(yaunTable);
}