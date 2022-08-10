window.addEventListener('DOMContentLoaded', function() {
    getImages();
});

function getImages() {
    async function callApi() {
        const res = await fetch('http://localhost:8081/images');
        const resObj = await res.json();
        const code = resObj["code"];
        console.log(code);
    }
    callApi();
}