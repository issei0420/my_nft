const uploadBox = document.getElementById('upload-box');
const fileInput = document.getElementById('file-input');

uploadBox.addEventListener('dragover', function(e) {
    e.stopPropagation();
    e.preventDefault();
    this.style.background = '#e1e7f0';   
}, false);

uploadBox.addEventListener('dragleave', function(e) {
    e.stopPropagation();
    e.preventDefault();
    this.style.background = '#f2f2f2';
}, false)

uploadBox.addEventListener('drop', function(e) {
    e.stopPropagation();
    e.preventDefault();
    this.style.background = '#f2f2f2';
    const files = e.dataTransfer.files;
    if (files.length > 1) return alert('アップロードできるファイルは1つだけです。');
    fileInput.files = files;
}, false);

fileInput.addEventListener('change', function(e) {
    console.log("ファイル名：" + e.target.files[0].name);
    document.getElementById('file-name').value = e.target.files[0].name;
})
