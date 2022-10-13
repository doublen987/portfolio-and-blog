function readURL(input) {
    return () => {
        if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            let thumbnailImg = document.getElementById("tag-thumbnail-image");
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


function initTagEditor() {

    window.onload = function() {
        var select = document.getElementById("chosen-tag");
        select.addEventListener("change", () => {
            if(select.value == "None") {
                document.getElementById("tag-name").value = "";
                document.getElementById("tag-description").value = "";
                document.getElementById("tag-thumbnail-image").src = "/content/no-image.png";
                document.getElementById("tag-thumbnail-name").value = "";
                document.getElementById("tah-thumbnailstretched").value = "false";
                var publish = document.getElementById("submit-publish");
                publish.classList.remove("submit-visible")
                publish.classList.add("submit-hidden")
                return
            }
            var content = document.getElementById("content-" + select.value).innerHTML;
            var title = document.getElementById("name-" + select.value).innerHTML;
            var description = document.getElementById("description-" + select.value).innerHTML;
            var thumbnail = document.getElementById("thumbnail-" + select.value).innerHTML;
            var thumbnailstretched = document.getElementById("thumbnailstretched-" + select.value).innerHTML;
            
            //console.log(decodeURI(stuff))
            console.log(he.decode(content))
            console.log(title);
            console.log(description);
            console.log(thumbnailstretched);
            document.getElementById("tag-name").value = title;
            document.getElementById("tag-description").value = description;
            document.getElementById("tag-thumbnail-name").value = thumbnail;
            document.getElementById("tag-thumbnailstretched").value = thumbnailstretched;
            var thumbnailImg = document.getElementById("tag-thumbnail-image");
            console.log("bla")
            thumbnailImg.src = "/content/images/" + thumbnail;
            thumbnailImg.width = 150;
            thumbnailImg.height = 200;
            thumbnailImg.onerror = onImageError;
            // var publish = document.getElementById("submit-publish");
            // publish.classList.remove("submit-hidden");
            // publish.classList.add("submit-visible");
        })
        
        let thumbnailInput = document.getElementById("tag-thumbnail");
        thumbnailInput.onchange = readURL(thumbnailInput);
    }
}