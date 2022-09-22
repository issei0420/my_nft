window.addEventListener("DOMContentLoaded", () => {
    const hiddenInputs = document.getElementsByClassName("image_id");
    let imageIds = [];
    for (elem of hiddenInputs) {
        imageIds.push(elem.value);
    }
    getPortionTotal(imageIds).then(total => {
        showData = shapeData(total);
        setTooltip(showData);
    })
}, false);

async function getPortionTotal(imageIds) {
    url = 'http://localhost:8080/portion'
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(imageIds)
    });
    return res.json();
}

function shapeData(total) {
    let showData = {}
    for (elem of total) {
        showData[elem["Id"]] = {}
    }
    for (elem of total) { 
        showData[elem["Id"]][elem["FileName"]] = elem["Count"]
    }
    return showData
}

function setTooltip (showData) {
    for (id in showData) {
        const td = document.getElementById(id);
        // ツールチップの準備
        td.setAttribute("data-bs-toggle", "tooltip");
        td.setAttribute("data-bs-html", "true");
        td.setAttribute("data-bs-placement", "right");
        text = createText(showData, id);
        td.setAttribute("data-bs-original-title", text)
    }
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
    })
}

function createText(showData, id) {
    totalData = showData[id]
    text = ""
    for (fileName in totalData) {
        text = text + `<p>${fileName}  ${totalData[fileName]}枚</p>`
    }
    return text;
}