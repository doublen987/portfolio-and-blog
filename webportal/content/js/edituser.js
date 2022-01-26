function readURL(input) {
    return () => {
        if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            let thumbnailImg = document.getElementById("user-thumbnail-image");
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

window.onload = function() {
    var select = document.getElementById("chosen-user");
    select.addEventListener("change", () => {
        if(select.value == "None") {
            editor.setContent(he.decode(""));
            document.getElementById("user-username").value = "";
            document.getElementById("user-thumbnail-image").src = "/content/no-image.png";
            document.getElementById("user-thumbnail-name").value = "";
            return
        }
        var username = document.getElementById("username-" + select.value).innerHTML;
        var description = document.getElementById("description-" + select.value).innerHTML;
        var thumbnail = document.getElementById("thumbnail-" + select.value).innerHTML;
        
        //console.log(decodeURI(stuff))
        console.log(username);
        console.log(description);
        document.getElementById("user-username").value = username;
        document.getElementById("user-description").value = description;
        document.getElementById("user-thumbnail-name").value = thumbnail;
        var thumbnailImg = document.getElementById("user-thumbnail-image");
        console.log("bla")
        thumbnailImg.src = "/content/images/" + thumbnail;
        thumbnailImg.width = 150;
        thumbnailImg.height = 200;
        thumbnailImg.onerror = onImageError;
    })
    
    let thumbnailInput = document.getElementById("user-thumbnail");
    thumbnailInput.onchange = readURL(thumbnailInput);
}


function initTinyMCE() {
    tinymce.init({
        selector:'#post-content',
        plugins: 'codesample image',
        toolbar: 'codesample',
        automatic_uploads: true,
        images_upload_url: '/content/images',
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
            var select = document.getElementById("chosen-user");
            select.addEventListener("change", () => {
                if(select.value == "None") {
                    editor.setContent(he.decode(""));
                    document.getElementById("user-username").value = "";
                    document.getElementById("user-password").value = "";
                    document.getElementById("user-thumbnail-image").src = "/content/no-image.png";
                    document.getElementById("user-thumbnail-name").value = "";
                    return
                }
                var username = document.getElementById("username-" + select.value).innerHTML;
                var password = document.getElementById("password-" + select.value).innerHTML;
                //var description = document.getElementById("description-" + select.value).innerHTML;
                var thumbnail = document.getElementById("thumbnail-" + select.value).innerHTML;
                
                //console.log(decodeURI(stuff))
                console.log(he.decode(content))
                console.log(username);
                console.log(password);
                console.log(description);
                console.log(published)
                console.log(publishtimestamp)
                editor.setContent(he.decode(content));
                document.getElementById("user-username").value = title;
                document.getElementById("user-password").value = description;
                document.getElementById("user-thumbnail-name").value = thumbnail;
                var thumbnailImg = document.getElementById("user-thumbnail-image");
                console.log("bla")
                thumbnailImg.src = "/content/images/" + thumbnail;
                thumbnailImg.width = 150;
                thumbnailImg.height = 200;
                thumbnailImg.onerror = onImageError;
            })
            
            let thumbnailInput = document.getElementById("user-thumbnail");
            thumbnailInput.onchange = readURL(thumbnailInput);
            

        }
    });
}