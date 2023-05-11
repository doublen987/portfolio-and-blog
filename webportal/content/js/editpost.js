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


function initTinyMCE() {
    tinymce.init({
        selector:'#post-content',
        plugins: [
            'advlist autolink lists link image charmap print preview anchor',
            'searchreplace visualblocks codesample fullscreen code',
            'insertdatetime media table paste codesample help wordcount tiny_mce_wiris'
        ],
        toolbar: 'undo redo | formatselect | ' +
        'bold italic backcolor | alignleft aligncenter ' +
        'alignright alignjustify | bullist numlist outdent indent | ' +
        'removeformat | help | image | codesample |' +
        'tiny_mce_wiris_formulaEditor | tiny_mce_wiris_formulaEditorChemistry | code',
        automatic_uploads: true,
        images_upload_url: '/content/images',
        image_dimensions: false,
        image_class_list: [
            {title: 'Responsive', value: 'img-responsive'}
        ],
        external_plugins: { tiny_mce_wiris: 'https://www.wiris.net/demo/plugins/tiny_mce/plugin.js' },
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
        },
        init_instance_callback : function(editor) {
            editor.setContent("<p>Hello world!</p>");
            var select = document.getElementById("chosen-post");
            select.addEventListener("change", () => {
                if(select.value == "None") {
                    editor.setContent(he.decode(""));
                    document.getElementById("post-title").value = "";
                    document.getElementById("post-description").value = "";
                    document.getElementById("post-thumbnail-image").src = "/content/no-image.png";
                    document.getElementById("post-thumbnail-name").value = "";
                    document.getElementById("post-thumbnailstretched").value = "false";
                    document.getElementById("post-publish-date").innerHTML = "";
                    document.getElementById("post-hidden").value = "False";
                    document.getElementById("post-published").value = "False";
                    document.getElementById("post-tags").innerHTML = "";
                    var publish = document.getElementById("submit-publish");
                    publish.classList.remove("submit-visible")
                    publish.classList.add("submit-hidden")
                    return
                }
                var content = document.getElementById("content-" + select.value).innerHTML;
                var title = document.getElementById("title-" + select.value).innerHTML;
                var description = document.getElementById("description-" + select.value).innerHTML;
                var thumbnail = document.getElementById("thumbnail-" + select.value).innerHTML;
                var thumbnailstretched = document.getElementById("thumbnailstretched-" + select.value).innerHTML;
                var publishtimestamp = document.getElementById("publishtimestamp-" + select.value).innerHTML;
                var lastedittimestamp = document.getElementById("lastedittimestamp-" + select.value).innerHTML;
                var hidden = document.getElementById('hidden-' + select.value).innerHTML;
                var published = document.getElementById('published-' + select.value).innerHTML;
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
                
                //console.log(decodeURI(stuff))
                console.log(he.decode(content))
                console.log(title);
                console.log(description);
                console.log(thumbnailstretched);
                console.log(published)
                console.log(publishtimestamp)
                editor.setContent(he.decode(content));
                document.getElementById("post-title").value = title;
                document.getElementById("post-description").value = description;
                document.getElementById("post-publish-date").innerHTML = publishtimestamp;
                document.getElementById("post-last-edit-date").innerHTML = lastedittimestamp;
                document.getElementById("post-thumbnail-name").value = thumbnail;
                document.getElementById("post-thumbnailstretched").value = thumbnailstretched;
                document.getElementById("post-hidden").value = hidden;
                document.getElementById("post-published").innerHTML = published;
                let postTags = document.getElementById("tag-section-tags")
                postTags.innerHTML = "";
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
                console.log("bla")
                thumbnailImg.src = "/content/images/" + thumbnail;
                thumbnailImg.width = 150;
                thumbnailImg.height = 200;
                thumbnailImg.onerror = onImageError;
                var publish = document.getElementById("submit-publish");
                publish.classList.remove("submit-hidden");
                publish.classList.add("submit-visible");
            })
            
            let thumbnailInput = document.getElementById("post-thumbnail");
            thumbnailInput.onchange = readURL(thumbnailInput);
            

        }
    });
}

export {
    initTinyMCE
}