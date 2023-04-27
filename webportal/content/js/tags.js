//let tagIndex = 0;

function initTags(id) {
    let container = document.getElementById(id)
    container.innerHTML = `
    <div class="choose-tag-modal" id="choose-tag-modal">
        
    </div>
    <div class="tags-section" id="tags-section">
        <div class="tag-section-tags" id="tag-section-tags">
        </div>
        <div class="tag-section-button-container">
            <button class="tag-section-button" id="tag-section-button">Add tag</button>
        </div>
    </div>`

    let addSectionItemButton = document.getElementById("tag-section-button")
    addSectionItemButton.addEventListener("click", showChooseTagModal(id))
    initTagModal()
}

function addTag(tag) {
    return (e) => {
        e.preventDefault()
        e.stopPropagation()
        let typeModal = document.getElementById("choose-tag-modal")
        typeModal.style.visibility = "hidden"
        let tags = document.getElementsByClassName("tag-section-tags")[0]

        let tagContainers = document.getElementsByClassName("tag-container")
        let tagIndex = tagContainers.length

        let id = Date.now()
        let container = document.createElement("div")
        container.classList.add("tag-container")
        container.id = ("tag-container-"+ id)
        container.innerHTML = `
            <div style="display:none;" class="stack-section-tag-ID" id="stack-section-tag-ID-${id}">${tag.ID}</div>
            <input style="display:none;" name="tag-${tagIndex}" value="${tag.ID}"></input>
            <img class="stack-section-tag-image" src="/content/images/${tag.thumbnail}" id="stack-section-tag-image-${id}"></img>
            <div class="tag-x" id="tag-x-${id}">
            <svg width="10" height="10" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><!--! Font Awesome Pro 6.2.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2022 Fonticons, Inc. --><path d="M376.6 84.5c11.3-13.6 9.5-33.8-4.1-45.1s-33.8-9.5-45.1 4.1L192 206 56.6 43.5C45.3 29.9 25.1 28.1 11.5 39.4S-3.9 70.9 7.4 84.5L150.3 256 7.4 427.5c-11.3 13.6-9.5 33.8 4.1 45.1s33.8 9.5 45.1-4.1L192 306 327.4 468.5c11.3 13.6 31.5 15.4 45.1 4.1s15.4-31.5 4.1-45.1L233.7 256 376.6 84.5z"/></svg>
            </div>
        `         
        tags.appendChild(container)

        let x = document.getElementById("tag-x-"+id)
        x.addEventListener("click", removeTag(id))

        //tagIndex++;

    }
}

function removeTag(id) {
    return () => {
        let tags = document.getElementsByClassName("tag-section-tags")[0]
        let tag = document.getElementById("tag-container-"+id)
        tags.removeChild(tag)
        tagIndex--;
    }
}

function initTagModal() {
    let tags = document.getElementsByClassName("tag-info")
    let typeModal = document.getElementById("choose-tag-modal")

    for(let i = 0; i < tags.length; i++) {
        let modaltag = document.createElement('div')
        modaltag.classList.add('modal-tag')
        modaltag.innerHTML = `
        <img src="/content/images/${tags[i].getElementsByClassName("tag-thumbnail")[0].innerHTML}"/>
        <div style="visibility:hidden;"  class="modal-tag-ID">${tags[i].getElementsByClassName("tag-ID")[0].innerHTML}</div>
        <div style="visibility:hidden;"  class="modal-tag-thumbnail">${tags[i].getElementsByClassName("tag-thumbnail")[0].innerHTML}</div>
        `
        typeModal.appendChild(modaltag)
    }

    let modaltags = typeModal.getElementsByClassName("modal-tag")
    console.log(modaltags.length)
    for(var i = 0; i < modaltags.length; i++) {
        console.log(i)
        let tagID = modaltags[i].getElementsByClassName("modal-tag-ID")[0]
        let tagThumbnail = modaltags[i].getElementsByClassName("modal-tag-thumbnail")[0]
        console.log({ 
            ID: tagID.innerHTML,
            thumbnail: tagThumbnail.innerHTML
        })
        modaltags[i].addEventListener("click", addTag({ 
            ID: tagID.innerHTML,
            thumbnail: tagThumbnail.innerHTML
        }), true)
        modaltags[i].addEventListener("click", () => console.log("bla"), true)
    }
}

function showChooseTagModal() {
    return (e) => {
        e.preventDefault()
        e.stopPropagation()
        let typeModal = document.getElementById("choose-tag-modal")
        typeModal.style.visibility = "visible"
        
    }
}

export {
    initTags,
    removeTag
}