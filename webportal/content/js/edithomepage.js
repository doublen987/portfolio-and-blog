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
  
function onClickSectionContainer(id) {
    return () => {

        let inputContainers = document.getElementsByClassName("input-container")
        for(let i = 0; i < inputContainers.length; i++) {
            let inputContainerID = inputContainers[i].id
            let id = inputContainerID.replace("input-container-","")
            console.log(id)
            //onClickSectionContainerX(id)()
        }

        console.log(id)
        let inputContianer = document.getElementById("input-container-"+id)
        
        let preview = inputContianer.getElementsByClassName("homepage-section-preview")[0]
        preview.style.display = "none"
        let editor = inputContianer.getElementsByClassName("homepage-section-editor")[0]
        editor.style.display = "unset"

        
        initTinyMCEForHompeageSection("section-content-input-"+id)
    }
}

function onClickSectionContainerX(id) {
    return (e) => {
        if(e) {
            e.stopPropagation()
        }
        console.log(id)

        let inputContianer = document.getElementById("input-container-"+id)
        let preview = inputContianer.getElementsByClassName("homepage-section-preview")[0]
        preview.style.display = "unset"
        let editor = inputContianer.getElementsByClassName("homepage-section-editor")[0]
        editor.style.display = "none"

        let headerinput = document.getElementById("section-header-input-"+id)
        let contentinput = document.getElementById("section-content-input-"+id)

        let headerpreview = document.getElementById("header-"+id)
        let contentpreview = document.getElementById("content-"+id)

        headerpreview.innerHTML = "Header "+headerinput.value
        contentpreview.innerHTML = "Content: "+tinymce.get("section-content-input-"+id).getContent()
    }
}

function onCLickAddSection(type) {
    
    return () => {
        let id = Date.now()
        console.log(type)
        switch(type) {
            case "text":
                let container = document.createElement("div")
                container.classList.add("input-container")
                container.classList.add("text-section-input-container")
                container.id = ("input-container-"+ id)

                let editorcontainer = document.createElement("div")
                editorcontainer.classList.add("homepage-section-editor")
                let containerX = document.createElement("div")
                containerX.classList.add("homepage-section-input-container-x")
                containerX.innerHTML = "<svg width=\"20\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 384 512\"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d=\"M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z\"/></svg>"
                editorcontainer.appendChild(containerX)
                let headercontainer = document.createElement("div")
                let headerlabel = document.createElement("label")
                headerlabel.classList.add("text-section-label-header")
                headerlabel.classList.add("editor-label")
                headerlabel.innerHTML="Header: "
                let header = document.createElement("input")
                header.classList.add("text-section-input-header")
                header.id = "section-header-input-"+id
                headercontainer.appendChild(headerlabel)
                headercontainer.appendChild(header)
                let contentcontainer = document.createElement("div")
                let contentlabel = document.createElement("label")
                contentlabel.classList.add("text-section-label-content")
                contentlabel.classList.add("editor-label")
                contentlabel.innerHTML="Content: "
                let content = document.createElement("textarea")
                content.classList.add("text-section-input-content")
                content.id = "section-content-input-"+id
                let updateButton = document.createElement("button")
                updateButton.classList.add("editor-update-button")
                contentcontainer.appendChild(contentlabel)
                contentcontainer.appendChild(content)
                editorcontainer.appendChild(headercontainer)
                editorcontainer.appendChild(contentcontainer)
                //contentcontainer.appendChild(updateButton)
                container.appendChild(editorcontainer)


                let previewcontainer = document.createElement("div")
                previewcontainer.classList.add("homepage-section-preview")
                let previewtitle = document.createElement("div")
                previewtitle.classList.add("homepage-section-preview-header")
                previewtitle.id = "header-"+id
                previewtitle.innerHTML = "Header: "
                let previewcontent = document.createElement("div")
                previewcontent.classList.add("homepage-section-preview-content")
                previewcontent.id = "content-"+id
                previewcontent.innerHTML = "Content: "
                previewcontainer.appendChild(previewtitle)
                previewcontainer.appendChild(previewcontent)
                container.appendChild(previewcontainer)

                containerX.addEventListener("click", onClickSectionContainerX(id))
                                
                previewcontainer.addEventListener("click", onClickSectionContainer(id))

                let form = document.getElementById("form1")
                form.appendChild(container)
            break;
            case "image": {
                let container = document.createElement("div")
                container.classList.add("input-container")
                container.classList.add("image-section-input-container")
                container.id = ("input-container-"+ id)
                container.innerHTML = `
                <div class="homepage-section-preview" id="homepage-section-preview-${id}">
                    <img class="homepage-section-image" id="homepage-section-image-${id}" src="">
                </div>
                <div class="homepage-section-editor" id="homepage-section-editor-${id}">
                    <div class="homepage-section-input-container-x">
                        <svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
                    </div>
                    <div>
                        <label for="homepage-section-editor-input" class="editor-label">Image: </label>
                        <input name="homepage-section-editor-input" type="file" accept="image/*" id="homepage-section-editor-input-${id}"></input>
                    </div>    
                </div>`


                container.getElementsByClassName("homepage-section-input-container-x")[0].addEventListener("click", onClickSectionContainerX(id))
                container.getElementsByClassName("homepage-section-preview")[0].addEventListener("click", onClickSectionContainer(id))
                let form = document.getElementById("form1")
                form.appendChild(container)
            }
            break;
            case "stack":
            break;
            default:
        }

        let typeModal = document.getElementById("homepage-section-type-modal")
        typeModal.style.display = "none"
    }
}

function showSectionTypeOptions() {
    let typeModal = document.getElementById("homepage-section-type-modal")
    typeModal.style.display = "unset"
}

function onSave() {
    let sections = []
    let sectionContainers = document.getElementsByClassName("homepage-section-input-container")
    for(let i = 0; i < sectionContainers.length; i++) {
        let currentsection = sectionContainers[i]
        if(currentsection.classList.contains("text-section-input-container")) {
            let header = currentsection.getElementsByClassName("text-section-input-header")[0]
            let content = currentsection.getElementsByClassName("text-section-input-content")[1]
            sections.push({
                type: "text",
                header: header,
                content: content
            })
        }
    }
}


window.onload = function() {
    let showTypeModalButton = document.getElementById("show-type-modal-button")
    showTypeModalButton.addEventListener("click", showSectionTypeOptions)

    let savebutton = document.getElementById("save-homepage-sections-button")
    savebutton.addEventListener("click", onSave)

    let types = [
        "text",
        "stack",
        "image"
    ]

    types.forEach(type => {
        let typeButton = document.getElementById("type-modal-button-"+type)
        typeButton.addEventListener("click", onCLickAddSection(type))
    })

    let inputContainers = document.getElementsByClassName("input-container")
    console.log(inputContainers)
    for(let i = 0; i < inputContainers.length; i++) {
        let container = inputContainers[i]
        let id = container.id.replace("input-container-","")
        container.addEventListener("click", onClickSectionContainer(id))
        let containersX = container.getElementsByClassName("homepage-section-input-container-x")[0]
        
        containersX.addEventListener("click", onClickSectionContainerX(id))
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


function initTinyMCEForHompeageSection(selector) {
    console.log("Initializing tinymce: "+selector)
    tinymce.init({
        selector: "#"+selector,
        plugins: 'codesample image',
        toolbar: 'codesample',
        automatic_uploads: true,
        images_upload_url: '/content/images',
        image_dimensions: false,
        image_class_list: [
            {title: 'Responsive', value: 'img-responsive'}
        ],
        relative_urls: false,
        file_picker_callback: function (cb, value, meta) {
            var input = document.createElement('input');
            input.setAttribute('type', 'file');
            input.setAttribute('accept', 'image/*');

            /*
            Note: In modern browsers input[type="file"] is functional without
            even adding it to the DOM, but that might not be the case in some older
            or quirky browsers like IE, so you might want to add it to the DOM
            just in case, and visually hide it. And do not forget do remove it
            once you do not need it anymore.
            */

            input.onchange = function () {
            var file = this.files[0];

            var reader = new FileReader();
            reader.onload = function () {
                /*
                Note: Now we need to register the blob in TinyMCEs image blob
                registry. In the next release this part hopefully won't be
                necessary, as we are looking to handle it internally.
                */
                var id = 'blobid' + (new Date()).getTime();
                var blobCache =  tinymce.activeEditor.editorUpload.blobCache;
                var base64 = reader.result.split(',')[1];
                var blobInfo = blobCache.create(id, file, base64);
                blobCache.add(blobInfo);

                console.log(blobInfo.blobUri())
                console.log(file.name)

                /* call the callback and populate the Title field with the file name */
                cb(blobInfo.blobUri(), { title: file.name });
            };
            reader.readAsDataURL(file);
            };

            input.click();
        },
        images_upload_handler: function(blobInfo, success, failure) {
            console.log("Uploading image")
            var xhr, formData;
            xhr = new XMLHttpRequest();
            xhr.withCredentials = false;
            xhr.open('POST', '/content/images')
            xhr.onload = function() {
                var json;
                if (xhr.status != 200) {
                    failure('HTTP Error: ' + xhr.status);
                    return;
                }
                json = JSON.parse(xhr.responseText);
                console.log(xhr.responseText)
                console.log(json)
                if (!json || typeof json.location != 'string') {
                    failure('Invalid JSON: ' + xhr.responseText);
                    return;
                }
                success("/content/images/" + json.location);
            };
            
            var jsonImage = JSON.stringify({
                filename: blobInfo.filename(),
                bytes: blobInfo.base64()
            })
            xhr.send(jsonImage)
        }
    });
}