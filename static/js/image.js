window.addEventListener('DOMContentLoaded', function() {
    showImages();
});

function showImages() {
    const divImage = document.getElementById("divImage")
    const filename = document.getElementById("fileName").innerHTML
    async function callApi() {
        const res = await fetch(`http://localhost:8080/images?filename=${filename}`);
        const resObj = await res.json();
        divImage.style.width = resObj["width"] + "px";
        for (const elem of resObj["images"]) {
            const imageElement = document.createElement("img")
            imageElement.src = "data:image/png;base64," + elem["code"]
            divImage.appendChild(imageElement);
        }
    }
    callApi();
}