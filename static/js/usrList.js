window.addEventListener("DOMContentLoaded", () => {
    const hiddenInputs = document.getElementsByClassName("image_id");
    let imageIds = [];
    for (elem of hiddenInputs) {
        imageIds.push(elem.value);
    }

    getPortionTotal(imageIds).then(total => {
        if (imageIds.length) {
            showData = shapeData(total);
            setTooltip(showData);
        }
        const sellerTable = document.getElementById("seller-table");
        if (sellerTable == undefined) {
            activateToolTip();
            return
        }
        getSellerImage().then(images => {
            showInfo = shapeInfo(images);
            setTooltipInfo(showInfo);
            activateToolTip();
        })
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
}

function createText(showData, id) {
    totalData = showData[id]
    text = ""
    for (fileName in totalData) {
        text = text + `<p>${fileName}  ${totalData[fileName]}枚</p>`
    }
    return text;
}

async function getSellerImage() {
    const res = await fetch(`http://localhost:8080/uploaded`)
    return res.json();
}

function shapeInfo(images) {
    showInfo = {}
    for (elem of images) {
        showInfo[elem["Id"]] = []
    }
    for (elem of images) {
        showInfo[elem["Id"]].push(elem["FileName"])
    }
    return showInfo
}

function activateToolTip() {
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
    });
}

function setTooltipInfo(showInfo) {
    for (id in showInfo){
        const td = document.getElementById(`i-${id}`);
        td.setAttribute("data-bs-toggle", "tooltip");
        td.setAttribute("data-bs-html", "true");
        td.setAttribute("data-bs-placement", "right");
        info = createInfo(showInfo, id)
        td.setAttribute("data-bs-original-title", info);
    }
}

function createInfo(showInfo, id) {
    info = ""
    for (filename of showInfo[id]) {
        info +=  `<p>${filename}</p>`
    }
    return info;
}