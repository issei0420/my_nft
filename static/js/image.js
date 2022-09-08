window.addEventListener('DOMContentLoaded', function() {
    showImages();
});

let portionInfo = {};
let flag = true;

const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
  return new bootstrap.Tooltip(tooltipTriggerEl)
})

function showImages() {
    const divImage = document.getElementById("divImage")
    const filename = document.getElementById("fileName").innerHTML
    async function callApi() {
        const res = await fetch(`http://3.114.104.27:8000/images?filename=${filename}`);
        const resObj = await res.json();
        divImage.style.width = resObj["width"] + "px";
        // 画像を表示
        for (const elem of resObj["images"]) {
            const imageElement = document.createElement("img")
            imageElement.src = "data:image/png;base64," + elem["code"]
            imageElement.setAttribute("id", elem["id"])
            divImage.appendChild(imageElement);
        }
        // portionInfoの準備
        for (const elem of resObj["portions"]) {
            portionInfo[elem["pId"]] = {
                "userId": elem["userId"], 
                "familyName": elem["familyName"],
                "firstName": elem["firstName"],
                "date": elem["date"].substring(0, elem["date"].indexOf(" ")),
            }
        }
        // 抽選部の設定
        for (const key in portionInfo) {
            const img = document.getElementById(key);
            // グレースケール化の関数を画像ロード時に実行
            img.onload = function() {
                if (!img.classList.contains("gray")) {
                    fillGray(img);
                    img.setAttribute("class", "gray")
                }
            }
            // ツールチップの準備
            img.setAttribute("data-bs-toggle", "tooltip");
            img.setAttribute("data-bs-html", "true");
            text = createText(key);
            img.setAttribute("data-bs-original-title", text);
        }
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
        const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
        })
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

function createText(id) {
    portion = portionInfo[id]
    const userId = portion["userId"]
    const name = portion["familyName"] + portion["firstName"]
    const date = portion["date"]
    const text = `<p>---- 所有者情報 ----</p><p>ユーザID：${userId}<br>姓名：${name}<br>抽選日：${date}</p>`
    return text;
}