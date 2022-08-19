window.addEventListener('DOMContentLoaded', function() {
    showImages();
});

let portionData = {};
let flag = true;

function showImages() {
    const divImage = document.getElementById("divImage")
    const filename = document.getElementById("fileName").innerHTML
    async function callApi() {
        const res = await fetch(`http://localhost:8080/images?filename=${filename}`);
        const resObj = await res.json();
        divImage.style.width = resObj["width"] + "px";
        // 画像を表示
        for (const elem of resObj["images"]) {
            const imageElement = document.createElement("img")
            imageElement.src = "data:image/png;base64," + elem["code"]
            imageElement.setAttribute("id", elem["id"])
            divImage.appendChild(imageElement);
        }
        // 抽選部位データをセット
        for (const elem of resObj["portions"]) {
            portionData[elem["pId"]] = {
                "userId": elem["userId"], 
                "familyName": elem["familyName"],
                "firstName": elem["firstName"],
                "data": elem["date"],
            }
        }
        // グレースケール化の関数を画像ロード時に実行
        for (const key in portionData) {
            const img = document.getElementById(key);
            img.onload = function() {
                if (!img.classList.contains("gray")) {
                    fillGray(img);
                    img.setAttribute("class", "gray")
                }
            }
        }
    }
    callApi();
}

function fillGray(img) {
    const width = img.width;
    const height = img.height;
    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;

    const context = canvas.getContext("2d");
    context.drawImage(img, 0, 0);

    const imgData = context.getImageData(0, 0, width, height);

    for (i = 0; i < height; i++) {
        for (j = 0; j < width; j++) {
            var pix = (i*width + j) * 4;
            var gray = 0.299 * imgData.data[pix] + 0.587 * imgData.data[pix+1] + 0.114 * imgData.data[pix+2];
            for (var k = 0; k < 3; k++) {
                imgData.data[pix+k] = gray;
           }
        }
    }
    context.putImageData(imgData, 0, 0)

    url = canvas.toDataURL();
    img.src = url;
}
