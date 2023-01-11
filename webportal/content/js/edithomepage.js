let sectionid = 0;
let pages = []
let readers = {}

function readURL(input, id) {
    return () => {
        if (input.files && input.files[0]) {
            var reader = new FileReader();
            reader.onload = function (e) {
                let thumbnailImg = document.getElementById("homepage-section-image-"+id);
                thumbnailImg.src = e.target.result;
                thumbnailImg.width = 150;
                thumbnailImg.height = 200;
                let container = document.getElementById('section-container-'+id)
                container.getElementsByClassName('homepage-image-filename')[0].innerHTML = input.files[0].name
            };
            reader.readAsDataURL(input.files[0]);
            readers[id] = reader
        }
    }
}
  
function onClickSectionContainer(id) {
    return () => {

        let inputContainers = document.getElementsByClassName("section-container")
        for(let i = 0; i < inputContainers.length; i++) {
            let inputContainerID = inputContainers[i].id
            let id = inputContainerID.replace("section-container-","")
            console.log(id)
            //onClickSectionContainerX(id)()
        }

        console.log(id)
        let inputContianer = document.getElementById("section-container-"+id)
        
        let preview = inputContianer.getElementsByClassName("homepage-section-preview")[0]
        preview.style.display = "none"
        let editor = inputContianer.getElementsByClassName("homepage-section-editor")[0]
        editor.style.display = "unset"

        let headerinput = document.getElementById("section-header-input-"+id)
        let contentinput = document.getElementById("section-content-input-"+id)

        

        initTinyMCEForHompeageSection("section-content-input-"+id, id)


        let headerpreview = document.getElementById("header-"+id)
        headerinput.value = headerpreview.innerText
    }
}

function onClickSectionContainerX(id) {
    return (e) => {
        if(e) {
            e.stopPropagation()
        }
        console.log(id)

        let inputContianer = document.getElementById("section-container-"+id)
        let preview = inputContianer.getElementsByClassName("homepage-section-preview")[0]
        preview.style.display = "unset"
        let editor = inputContianer.getElementsByClassName("homepage-section-editor")[0]
        editor.style.display = "none"

        let headerinput = document.getElementById("section-header-input-"+id)
        let contentinput = document.getElementById("section-content-input-"+id)

        let headerpreview = document.getElementById("header-"+id)
        let contentpreview = document.getElementById("content-"+id)

        headerpreview.innerHTML = headerinput.value
        contentpreview.innerHTML = tinymce.get("section-content-input-"+id).getContent()
    }
}

//Called when you click the modal button for a given section
function onCLickAddSection(section) {
    
    return () => {
        let id = Date.now()
        console.log(section)
        switch(section.type) {
            case "text":
                let container = document.createElement("div")
                container.classList.add("section-container")
                container.classList.add("text-section-section-container")
                container.id = ("section-container-"+ id)

                let blacontainerX = document.createElement("div")
                blacontainerX.id="remove-section-button-"+id
                blacontainerX.classList.add("homepage-section-container-x")
                blacontainerX.innerHTML = "<svg width=\"20\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 384 512\"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d=\"M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z\"/></svg>"
                container.appendChild(blacontainerX)

                let editorcontainer = document.createElement("div")
                editorcontainer.classList.add("homepage-section-editor")
                let containerX = document.createElement("div")
                containerX.classList.add("homepage-section-section-container-x")
                containerX.innerHTML = "<svg width=\"20\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 384 512\"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d=\"M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z\"/></svg>"
                editorcontainer.appendChild(containerX)
                let headercontainer = document.createElement("div")
                let headerlabel = document.createElement("label")
                headerlabel.classList.add("text-section-label-header")
                headerlabel.classList.add("editor-label")
                headerlabel.innerHTML= "Header: "
                let header = document.createElement("input")
                header.classList.add("text-section-input-header")
                header.id = "section-header-input-"+id
                headercontainer.appendChild(headerlabel)
                headercontainer.appendChild(header)
                let contentcontainer = document.createElement("div")
                let contentlabel = document.createElement("label")
                contentlabel.classList.add("text-section-label-content")
                contentlabel.classList.add("editor-label")
                contentlabel.innerHTML= "Content: "
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
                previewtitle.innerHTML = section.header? section.header : "Header"
                let previewcontent = document.createElement("div")
                previewcontent.classList.add("homepage-section-preview-content")
                previewcontent.id = "content-"+id
                previewcontent.innerHTML = section.content? section.content : "Content"
                previewcontainer.appendChild(previewtitle)
                previewcontainer.appendChild(previewcontent)
                container.appendChild(previewcontainer)

                containerX.addEventListener("click", onClickSectionContainerX(id))
                blacontainerX.addEventListener("click", onClickRemoveSection(id))                

                previewcontainer.addEventListener("click", onClickSectionContainer(id))

                let form = document.getElementById("sections-container")
                form.appendChild(container)
            break;
            case "image": {
                let sectionFilename = section.filename? section.filename : ""
                let container = document.createElement("div")
                container.classList.add("section-container")
                container.classList.add("image-section-section-container")
                container.id = ("section-container-"+ id)
                container.innerHTML = `
                <div class="homepage-section-preview" id="homepage-section-preview-${id}">
                    <div class="homepage-section-section-container-x">
                        <svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
                    </div>    
                    <img class="homepage-section-image" id="homepage-section-image-${id}" src="/content/images/${section.filename}">
                    
                    <div>
                        <label for="homepage-section-editor-input" class="editor-label">Image: </label>
                        <input name="homepage-section-editor-input" type="file" accept="image/*" id="homepage-section-editor-input-${id}"></input>
                        <div class="homepage-image-section-id" >${id}</div>
                        <div class="homepage-image-filename" >${sectionFilename}</div>
                    </div>    
                </div>`


                // container.getElementsByClassName("homepage-section-section-container-x")[0].addEventListener("click", onClickSectionContainerX(id))
                // container.getElementsByClassName("homepage-section-preview")[0].addEventListener("click", onClickSectionContainer(id))
                let form = document.getElementById("sections-container")
                form.appendChild(container)


                var thumbnailImg = document.getElementById(`homepage-section-image-${id}`);
                console.log("bla")
                thumbnailImg.width = 150;
                thumbnailImg.height = 200;
                thumbnailImg.onerror = onImageError;

                let thumbnailInput = document.getElementById(`homepage-section-editor-input-${id}`);
                thumbnailInput.onchange = readURL(thumbnailInput, id);
            }
            break;
            case "stack": {
                let container = document.createElement("div")
                container.classList.add("section-container")
                container.classList.add("stack-section-section-container")
                container.id = ("section-container-"+ id)
                container.innerHTML = `
                <div id="remove-section-button-${id}">
                    <svg width=\"20\" height=\"20\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 384 512\"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d=\"M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z\"/></svg>
                </div>
                <div class="tags-section" id="tags-section-${id}">
                    
                </div>
                <div>
                    <button id="add-tag-section-button">Add tag section</button>
                </div>`

                if (!section.tagssections) {
                    section.tagssections = []
                }
                
                let form = document.getElementById("sections-container")
                form.appendChild(container)

                //container.getElementsByClassName("homepage-section-section-container-x")[0].addEventListener("click", onClickSectionContainerX(id))
                document.getElementById("add-tag-section-button").addEventListener("click", addTagSection(id))

                document.getElementById("remove-section-button-"+id).addEventListener("click", onClickRemoveSection(id))    

                console.log(section.tagssections)
                for(let i = 0; i < section.tagssections.length; i++) {

                    let tagssection = convertMapToSection(section.tagssections[i])
                    console.log(tagssection)
                    let tagSectionID = addTagSection(id, tagssection.name)()
                    for(let j = 0; j <  tagssection.tags.length; j++) {
                        console.log(convertMapToSection(tagssection.tags[j]))
                        sectionid = tagSectionID
                        addTagToTagsSectionItem(convertMapToSection(tagssection.tags[j]))()
                    } 
                }
            }
            break;
            case "3dmodel": {
                let sectionFilename = section.filename? section.filename : ""
                let container = document.createElement("div")
                container.classList.add("section-container")
                container.classList.add("3dmodel-section-section-container")
                container.id = ("section-container-"+ id)
                container.innerHTML = `
                <div class="homepage-section-preview" id="homepage-section-preview-${id}">

                    <div id="homepage-section-3dmodel-${id}" ></div>

                </div>
                <div>
                    <label for="homepage-section-editor-input" class="editor-label">Model: </label>
                    <input name="homepage-section-editor-input" type="file" id="homepage-section-editor-input-${id}"></input>
                    <div class="homepage-3dmodel-section-id" >${id}</div>
                    <div class="homepage-3dmodel-filename" >${sectionFilename}</div>
                </div>`


                // container.getElementsByClassName("homepage-section-section-container-x")[0].addEventListener("click", onClickSectionContainerX(id))
                // container.getElementsByClassName("homepage-section-preview")[0].addEventListener("click", onClickSectionContainer(id))

                let form = document.getElementById("sections-container")
                form.appendChild(container)
                console.log(section)
                if(section.filename) {
                    let stl_viewer = init3dViewer(`homepage-section-3dmodel-${id}`, '/content/images/'+section.filename)
                    stl_viewer.add_model({id:1, filename: '/content/images/'+section.filename})
                }

                document.getElementById(`homepage-section-editor-input-${id}`).addEventListener("change", function() {
                    const file = this.files[0];
                    ///const viewer = document.getElementById(`homepage-section-3dmodel-${id}`);
                  
                    if (!file) return;
                  
                    console.log(file)

                    const reader = new FileReader();
                    reader.addEventListener('load', async function() {
                    //   const texture = await viewer.createTexture(reader.result);

                       console.log(reader.result)
                       let stl_viewer = init3dViewer(`homepage-section-3dmodel-${id}`)
                       stl_viewer.add_model({local_file:file});
                       readers[id] = reader
                       container.getElementsByClassName('homepage-3dmodel-filename')[0].innerHTML = file.name
                    });
                    reader.addEventListener('progress', () => {
                        console.log("bla")
                    })
                    reader.addEventListener('error', () => {
                        console.log("error")
                    })
                    reader.readAsDataURL(file);
                })


                
            }
            break;
            
            default:
        }

        let typeModal = document.getElementById("homepage-section-type-modal")
        typeModal.style.display = "none"
    }
}

function onExitSectionNameInput(id) {
    return () => {
        let name = document.getElementById("tag-section-name-"+id)
        name.style.display = "block"
        let inputContainer = document.getElementById("tag-section-section-container-"+id)
        inputContainer.style.display = "none"
        let input = document.getElementById("tag-section-input-"+id)
        name.innerHTML = input.value
    }
}

function addTagSection(sectionid, tsn) {
    return () => {
        let tagsSection = document.getElementById("tags-section-"+sectionid)
        let id = Date.now()
        let container = document.createElement("div")
        container.classList.add("tag-group-container")
        container.id = ("tag-group-container-"+ id)
        container.innerHTML = `
        <div class="tag-group-container-x" id="tag-group-container-x-${id}">
            <svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
        </div>
        <div class="tag-section-name-container">
            <div class="tags-section-header-preview" id="tag-section-name-${id}">${tsn}</div>
            <div style="display:none;" class="tag-section-section-container" id="tag-section-section-container-${id}">
                <div class="tag-section-input-x" id="tag-section-input-x-${id}">
                    <svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
                </div>
                <input class="tag-section-input" id="tag-section-input-${id}"></input>
                
            </div>
        </div>
        <div class="tag-section-tags" id="tag-section-tags-${id}">
        </div>
        <div class="tag-section-button-container">
            <button class="tag-section-button" id="tag-section-button-${id}">Add tag</button>
        </div>
        `                   
        tagsSection.appendChild(container)

        let tagSectionName = document.getElementById("tag-section-name-"+id)
        tagSectionName.addEventListener("click", onClickTagSectionName(id))

        let tagSectionInputX = document.getElementById("tag-section-input-x-"+id)
        tagSectionInputX.addEventListener("click", onExitSectionNameInput(id))

        let addSectionItemButton = document.getElementById("tag-section-button-"+id)
        addSectionItemButton.addEventListener("click", showChooseTagModal(id))

        let tagGroupX = document.getElementById("tag-group-container-x-"+id)
        tagGroupX.addEventListener("click", removeTagGroup(id))

        return id;
    }
}

function onChangeTagSectionNameInput(id) {
    return (e)=> {
        let newname = e.target.value
        let name = document.getElementById("tag-section-name-"+id)
        name.innerHTML = newname
    }
}

function onClickTagSectionName(id) {
    return () => {
        let name = document.getElementById("tag-section-name-"+id)
        name.style.display = "none"
        let input = document.getElementById("tag-section-section-container-"+id)
        input.style.display = "unset"
    }
}

function addTagToTagsSectionItem(tag) {
    return (e) => {
        if(e) {
            e.preventDefault()
            e.stopPropagation()
        }
        
        let sectionitemid = sectionid
        let typeModal = document.getElementById("choose-tag-modal")
        typeModal.style.visibility = "hidden"
        let tags = document.getElementById("tag-section-tags-"+sectionitemid)

        console.log("tag-section-tags-"+sectionitemid)

        let id = Date.now()
        let container = document.createElement("div")
        container.classList.add("tag-container")
        container.id = ("tag-container-"+ id)
        container.innerHTML = `
            <div style="display:none;" class="stack-section-tag-ID" id="stack-section-tag-ID-${id}">${tag.ID}</div>
            <div style="display:none;" class="stack-section-tag-name" id="stack-section-tag-name-${id}">${tag.name}</div>
            <div style="display:none;" class="stack-section-tag-thumbnail" id="stack-section-tag-thumbnail-${id}">${tag.thumbnail}</div>
            <img class="stack-section-tag-image" src="/content/images/${tag.thumbnail}" id="stack-section-tag-image-${id}"></img>
            <div class="tag-x" id="tag-x-${id}">
            <svg width="10" height="10" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
            </div>
        `         
        tags.appendChild(container)

        let x = document.getElementById("tag-x-"+id)
        x.addEventListener("click", removeTagFromTagsSectionItem(sectionitemid, id))

        

    }
}

function removeTagFromTagsSectionItem(sectionitemid, id) {
    return () => {
        let tags = document.getElementById("tag-section-tags-"+sectionitemid)
        let tag = document.getElementById("tag-container-"+id)
        tags.removeChild(tag)
    }
}

function initTagModal() {
    tags = document.getElementsByClassName("tag-info")
    let typeModal = document.createElement('div')
    typeModal.innerHTML = `<div class="choose-tag-modal" id="choose-tag-modal"></div>`
    document.getElementById("form1").appendChild(typeModal)

    for(let i = 0; i < tags.length; i++) {
        let modaltag = document.createElement('div')
        modaltag.classList.add('modal-tag')
        modaltag.innerHTML = `
        <img src="/content/images/${tags[i].getElementsByClassName("tag-thumbnail")[0].innerHTML}"/>
        <div style="visibility:hidden;" class="modal-tag-ID">${tags[i].getElementsByClassName("tag-ID")[0].innerHTML}</div>
        <div style="visibility:hidden;"   class="modal-tag-thumbnail">${tags[i].getElementsByClassName("tag-thumbnail")[0].innerHTML}</div>
        <div style="visibility:hidden;"   class="modal-tag-name">${tags[i].getElementsByClassName("tag-name")[0].innerHTML}</div>
        `
        console.log(modaltag)
        document.getElementById("choose-tag-modal").appendChild(modaltag)
    }

    modaltags = typeModal.getElementsByClassName("modal-tag")
    console.log(modaltags.length)
    for(var i = 0; i < modaltags.length; i++) {
        console.log(i)
        let tagID = modaltags[i].getElementsByClassName("modal-tag-ID")[0]
        let tagThumbnail = modaltags[i].getElementsByClassName("modal-tag-thumbnail")[0]
        let tagName = modaltags[i].getElementsByClassName("modal-tag-name")[0]
        modaltags[i].addEventListener("click", addTagToTagsSectionItem({ 
            ID: tagID.innerHTML,
            thumbnail: tagThumbnail.innerHTML,
            name: tagName.innerHTML
        }, sectionid), true)
        modaltags[i].addEventListener("click", () => console.log("bla"), true)
    }
}

function showChooseTagModal(id) {
    return () => {
        sectionid = id
        console.log(sectionid)
        let typeModal = document.getElementById("choose-tag-modal")
        typeModal.style.visibility = "visible"
        
    }
}

function showSectionTypeModal() {
    let typeModal = document.getElementById("homepage-section-type-modal")
    typeModal.style.display = "unset"
}

function addTextSectionToSaveList(sections, currentsection) {
    let header = currentsection.getElementsByClassName("homepage-section-preview-header")[0].innerHTML
    let content = currentsection.getElementsByClassName("homepage-section-preview-content")[0].innerHTML
    sections.push({
        type: "text",
        header:header,
        content:  content
    })
}

function addStackSectionToSaveList(sections, currentsection) {
    
    let tagGroupContainer = currentsection.getElementsByClassName("tag-group-container")

    let header = ""
    let tagssections = [

    ]

    for(let i = 0; i < tagGroupContainer.length; i++) {
        header = tagGroupContainer[i].getElementsByClassName("tags-section-header-preview")[0]
        let tagsHTML = tagGroupContainer[i].getElementsByClassName("tag-container")
        console.log(tagsHTML)
        let tags = []
        for(let i = 0; i < tagsHTML.length; i++) {
            let ID = tagsHTML[i].getElementsByClassName("stack-section-tag-ID")[0].innerHTML
            let name = tagsHTML[i].getElementsByClassName("stack-section-tag-name")[0].innerHTML
            let thumbnail = tagsHTML[i].getElementsByClassName("stack-section-tag-thumbnail")[0].innerHTML
            tags.push({
                ID: ID,
                name: name,
                thumbnail: thumbnail
            })
        }
        tagssections.push({
            name: header.innerHTML,
            tags: tags
        })

    }

    console.log(tagssections)
    sections.push({
        type: "stack",
        name: header.innerHTML,
        tagssections: tagssections
    })
}

function addImageSectionToSaveList(sections, currentsection) {
    let filename = currentsection.getElementsByClassName("homepage-image-filename")[0].innerHTML
    let id = currentsection.getElementsByClassName("homepage-image-section-id")[0].innerHTML
    sections.push({
        type: "image",
        filename:filename,
        bytes: readers[id] && readers[id].result? readers[id].result.replace(new RegExp("data:image/(png|jpg|jpeg);base64,"),'') : []
    })
}

function add3DModelSectionToSaveList(sections, currentsection) {
    let filename = currentsection.getElementsByClassName("homepage-3dmodel-filename")[0].innerHTML
    let id = currentsection.getElementsByClassName("homepage-3dmodel-section-id")[0].innerHTML
    sections.push({
        type: "3dmodel",
        filename: filename,
        bytes: readers[id] && readers[id].result? readers[id].result.replace(new RegExp("data:application/x-tgif;base64,"),'') : []
    })
}

function onClickRemoveSection(id) {
    return () => {
        document.getElementById("section-container-"+id).remove()
        readers[id] = undefined
    }
}

function removeTagGroup(id) {
    return () => {
        document.getElementById("tag-group-container-"+ id).remove()
    }
}

function onSave() {
    let sections = []
    let sectionContainers = document.getElementsByClassName("section-container")
    for(let i = 0; i < sectionContainers.length; i++) {
        let currentsection = sectionContainers[i]
        if(currentsection.classList.contains("text-section-section-container")) {
            addTextSectionToSaveList(sections, currentsection)
        }
        if(currentsection.classList.contains("stack-section-section-container")) {
            addStackSectionToSaveList(sections, currentsection)
        }
        if(currentsection.classList.contains("image-section-section-container")) {
            addImageSectionToSaveList(sections, currentsection)
        }
        if(currentsection.classList.contains("3dmodel-section-section-container")) {
            add3DModelSectionToSaveList(sections, currentsection)
        }

    }
    console.log(sections)

    let pageID = document.getElementById("page-select").value;
    let pageName = document.getElementById("page-name-input").value;
    let pageHomepageCheckmark = document.getElementById("page-homepage-checkmark").checked

    async function savePage() {
        const rawResponse = await fetch('/homepage/edit', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: pageID,
                homepage: pageHomepageCheckmark,
                name: pageName,
                sections: sections
            })    
        });
        const content = await rawResponse;


        console.log(rawResponse.ok);
        if(rawResponse.ok) {
        }
    }

    savePage();
}

function onDelete() {
    let pageID = document.getElementById("page-select").value;

    async function deletePage() {
        const rawResponse = await fetch('/homepage/edit', {
            method: 'DELETE',
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: pageID
            })    
        });
        const content = await rawResponse;


        console.log(rawResponse.ok);
        if(rawResponse.ok) {
        }
    }

    deletePage();
}

function getPages() {
    fetch('/pages', {
    method: 'GET',
    credentials: 'include',
    
    })
    .then(response => response.json())
    .then(data => {
        console.log(data)
        pages = data
    });
    

    let pageSelection = document.getElementById("page-select")
    pageSelection.addEventListener("change", onChangeSelectedPage)
    
}

function convertMapToSection(map) {
    let sectionObject = {}
    for(let i = 0; i < map.length; i++) {
        sectionObject[map[i].Key] = map[i].Value
    }
    return sectionObject
}

function showPage(page) {
    
    console.log(page.sections)
    for(let i = 0; i < page.sections.length; i++) {
        onCLickAddSection(convertMapToSection(page.sections[i]))()
    }
}

function onChangeSelectedPage(e) {
    let sections  = document.querySelectorAll(".section-container")
    let seclen = sections.length
    sections.forEach(section => {
        console.log(sections)
        section.remove()
    })
    let selection = e.target
    console.log(selection)
    let selectedPageID = selection.value
    for(let i = 0; i < pages.length; i++) {
        if(pages[i].ID == selectedPageID) {
            
            showPage(pages[i])
            let pageNameInput = document.getElementById("page-name-input")
            pageNameInput.value = pages[i].name

            let pageHomepageCheckmark = document.getElementById("page-homepage-checkmark")
            pageHomepageCheckmark.checked = pages[i].homepage
        }
    }

    
}

window.onload = function() {
    let showTypeModalButton = document.getElementById("show-type-modal-button")
    showTypeModalButton.addEventListener("click", showSectionTypeModal)

    let savebutton = document.getElementById("save-homepage-sections-button")
    savebutton.addEventListener("click", onSave)

    let deletebutton = document.getElementById("delete-homepage-sections-button")
    deletebutton.addEventListener("click", onDelete)

    let types = [
        "text",
        "stack",
        "image",
        "3dmodel"
    ]

    types.forEach(type => {
        let typeButton = document.getElementById("type-modal-button-"+type)
        typeButton.addEventListener("click", onCLickAddSection({type: type}))
    })

    let inputContainers = document.getElementsByClassName("section-container")
    console.log(inputContainers)
    for(let i = 0; i < inputContainers.length; i++) {
        let container = inputContainers[i]
        let id = container.id.replace("section-container-","")
        container.addEventListener("click", onClickSectionContainer(id))
        let containersX = container.getElementsByClassName("homepage-section-section-container-x")[0]
        
        containersX.addEventListener("click", onClickSectionContainerX(id))
    }

    initTagModal()
    getPages()

    var el = document.getElementById('sections-container');
    var sortable = Sortable.create(el);

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


function initTinyMCEForHompeageSection(selector, id) {
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
        },
        init_instance_callback : function(editor) {
            editor.setContent("<p>Hello world!</p>");
            let contentpreview = document.getElementById("content-"+id)
            editor.setContent(contentpreview.innerHTML)

        }
    });
}