function readURL(input) {
    return () => {
        if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            let thumbnailImg = document.getElementById("post-thumbnail-image");
            thumbnailImg.src = e.target.result;
            thumbnailImg.width = 150;
            thumbnailImg.height = 200;
        };
        reader.readAsDataURL(input.files[0]);
        }
    }
}

function onImageError() {
    //this.onerror=null;
    //if (this.src != '/content/no-image.png') {
        console.log("error!");
        console.log(this);
        this.src = '/content/no-image.png';
        this.width = 150;
        this.height = 200;
    //}
}

let thumbnailInput = document.getElementById("post-thumbnail");
thumbnailInput.onchange = readURL(thumbnailInput);

var select = document.getElementById("chosen-post");
select.addEventListener("change", () => {
    if(select.value == "None") {
        document.getElementById("post-title").value = "";
        document.getElementById("post-description").value = "";
        document.getElementById("post-link").value = "";
        document.getElementById("post-thumbnailstretched").value = "false";
        document.getElementById("post-thumbnail-name").value = "";
        return
    }
    var description = document.getElementById("description-" + select.value).innerHTML;
    var title = document.getElementById("title-" + select.value).innerHTML;
    var link = document.getElementById("link-" + select.value).innerHTML;
    var thumbnailstretched = document.getElementById("thumbnailstretched-" + select.value).innerHTML;
    var thumbnail = document.getElementById("thumbnail-" + select.value).innerHTML;
    console.log(thumbnailstretched)
    document.getElementById("post-description").value = description;
    document.getElementById("post-title").value = title;
    document.getElementById("post-link").value = link;
    document.getElementById("post-thumbnailstretched").value = thumbnailstretched;
    document.getElementById("post-thumbnail-name").value = thumbnail;
})