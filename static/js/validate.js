document.addEventListener('DOMContentLoaded', () => {
    const theForm = document.getElementById("form")
    const requiredElems = document.querySelectorAll('.required');
    const errorClassName = 'error';

    // 未入力の欄に警告エラーを表示
    const createError = (elem, errorMessage) => {
        const errorSpan = document.createElement('span');
        errorSpan.classList.add(errorClassName);
        errorSpan.setAttribute('aria-live', 'polite');
        errorSpan.textContent = errorMessage;
        errorSpan.style.color = "#dc3545";
        elem.parentNode.appendChild(errorSpan);
    }    

    theForm.addEventListener("submit", (e) => {
    //初期化
      const errorElems = theForm.querySelectorAll('.' + errorClassName);
      errorElems.forEach( (elem) => {
        elem.remove(); 
      });

    //.required を指定した要素を検証
    requiredElems.forEach( (elem) => {
        const elemValue = elem.value.trim(); e
        if(elemValue.length === 0) {
            createError(elem, '入力は必須です');
            e.preventDefault();
        }
        });
    })
});

var lotteryUnits = 0;
const selectUserType = document.getElementById("selectUserType")
const inputLotteryUnits = document.getElementById("inputLotteryUnits")
selectUserType.addEventListener("change", (e) => {
    if (selectUserType.value == "sellers") {
      inputLotteryUnits.disabled = true;
      inputLotteryUnits.value = "";
    } else {
      inputLotteryUnits.disabled = false;
      inputLotteryUnits.value = "0";
    }
    e.preventDefault();
});