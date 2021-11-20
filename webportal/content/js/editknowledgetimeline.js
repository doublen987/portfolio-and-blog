function InitializeEventEditor() {
    var select = document.getElementById("chosen-event");
    select.addEventListener("change", () => {
        if(select.value == "None") {
            editor.setContent(he.decode(""));
            document.getElementById("event-title").value = "";
            return
        }
        var content = document.getElementById("description-" + select.value).innerHTML;
        var title = document.getElementById("title-" + select.value).innerHTML;
        //console.log(stuff)
        //console.log(decodeURI(stuff))
        console.log(he.decode(content))
        console.log(title)
        editor.setContent(he.decode(content));
        document.getElementById("event-description").value = title;
        document.getElementById("event-title").value = title;
    })
}