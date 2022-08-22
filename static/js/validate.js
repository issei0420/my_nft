document.addEventListener('DOMContentLoaded', () => {
    const theForm = document.getElementById("form")
    const requiredElems = document.querySelectorAll('.required');
    const errorClassName = 'error';

    const createError = (elem, errorMessage) => {
        //span 要素を生成
        const errorSpan = document.createElement('span');
        //エラー用のクラスを追加（設定）
        errorSpan.classList.add(errorClassName);
        //aria-live 属性を設定
        errorSpan.setAttribute('aria-live', 'polite');
        //引数に指定されたエラーメッセージを設定
        errorSpan.textContent = errorMessage;
        // テキストを赤色にするcssを設定
        errorSpan.style.color = "#dc3545";
        //elem の親要素の子要素として追加
        elem.parentNode.appendChild(errorSpan);
    }    

    theForm.addEventListener("submit", (e) => {
    //エラーを表示する要素を全て取得して削除（初期化）
      const errorElems = theForm.querySelectorAll('.' + errorClassName);
      errorElems.forEach( (elem) => {
        elem.remove(); 
      });

    //.required を指定した要素を検証
    requiredElems.forEach( (elem) => {
        //値（value プロパティ）の前後の空白文字を削除
        const elemValue = elem.value.trim(); 
        //値が空の場合はエラーを表示してフォームの送信を中止
        if(elemValue.length === 0) {
            createError(elem, '入力は必須です');
            e.preventDefault();
        }
        });
    })
});

// const selectUserType = document.getElementById("selectUserType")
// selectUserType.addEventListener("change", () => {
//     if selectUserType.value == 
// })