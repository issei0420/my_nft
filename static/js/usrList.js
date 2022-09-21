window.addEventListener("DOMContentLoaded", () => {
    const hiddenInputs = document.getElementsByClassName("image_id");
    let imageIds = [];
    for (elem of hiddenInputs) {
        imageIds.push(elem.value);
    }
    getPortionTotal(imageIds).then(total => {
        shapeData(total)
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
        showData[elem["Id"]] = []
    }
    for (elem of total) { 
        showData[elem["Id"]].push([elem["FileName"], elem["Count"]])
    }
}

