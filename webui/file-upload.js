
// grab the form data and upload the file to RESTCONF rpc
function upload() {
    const bookName = document.getElementById("bookName").value;
    const pdf = document.getElementById("pdf").files[0];
    const form = new FormData();
    form.append("bookName", bookName); // YANG: leaf bookName { type string; }
    form.append("pdf", pdf);           // YANG: anydata pdf;
    fetch("/restconf/data/file-upload:bookReport", {
      method: "POST",
      body: form
    })
    .then(resp => resp.json())
    .then((data) => {
        window.alert(`success. file size is ${data.fileSize}`);
    });
}