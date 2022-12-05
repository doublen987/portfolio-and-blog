import {removeTag} from "/content/js/tags.js"

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

window.onload = () => {

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
    var tagsElement = document.getElementById('tags-' + select.value);
    var tagElements = tagsElement.getElementsByClassName("tag")
    let tags = []
    for(let i = 0; i < tagElements.length; i++) {
        let tagid = tagElements[i].getElementsByClassName("tag-ID")[0].innerHTML
        let tagthumbnail = tagElements[i].getElementsByClassName("tag-thumbnail")[0].innerHTML
        tags.push({
            ID: tagid,
            thumbnail: tagthumbnail
        })
    }

    console.log(thumbnailstretched)
    document.getElementById("post-description").value = description;
    document.getElementById("post-title").value = title;
    document.getElementById("post-link").value = link;
    document.getElementById("post-thumbnailstretched").value = thumbnailstretched;
    document.getElementById("post-thumbnail-name").value = thumbnail;

    let postTags = document.getElementById("tag-section-tags")
    postTags.innerHTML = "";
    console.log(tags)
    for(let i = 0; i < tags.length; i++) {
        let id = Date.now()
        let container = document.createElement("div")
        container.classList.add("tag-container")
        container.id = ("tag-container-"+ id)
        container.innerHTML = `
            <div style="display:none;" class="stack-section-tag-ID" id="stack-section-tag-ID-${id}">${tags[i].ID}</div>
            <input style="display:none;" name="tag-${i}" value="${tags[i].ID}"></input>
            <img class="stack-section-tag-image" src="/content/images/${tags[i].thumbnail}" id="stack-section-tag-image-${id}"></img>
            <div class="tag-x" id="tag-x-${id}">
            <svg width="10" height="10" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
            </div>
        `         
        postTags.appendChild(container)

        let x = document.getElementById("tag-x-"+id)
        x.addEventListener("click", removeTag(id))

    }

    var thumbnailImg = document.getElementById("post-thumbnail-image");
    console.log(thumbnail)
    thumbnailImg.src = "/content/images/" + thumbnail;
    thumbnailImg.width = 150;
    thumbnailImg.height = 200;
    thumbnailImg.onerror = onImageError;
})
}