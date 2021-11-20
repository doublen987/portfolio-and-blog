function manipulateElement(reveal, id) {
    var timelineEvent = document.getElementById("timeline-event-" + id);
    if(timelineEvent) {
        if(reveal) {
            timelineEvent.classList.remove("event-not-visible")
            timelineEvent.classList.add("event-visible")
        } else {
            timelineEvent.classList.remove("event-visible")
            timelineEvent.classList.add("event-not-visible")
        }
    }
}

var previousFirst

var onScroll = function() {
    //var cont = document.getElementsByClassName("timeline-container")
    var lastID = Math.trunc(window.scrollY / 167);
    var length = Math.trunc(window.innerHeight / 167);
    
    for(let i = 0; i < lastID; i++) {
        manipulateElement(false, i)
    }

    for(let v = lastID; v < lastID + length; v++) {
        manipulateElement(true, v)
    }

    for(let j = lastID + length; j < 100; j++) {
        manipulateElement(false, j)
    }

    console.log("first ID" + Math.trunc(lastID) + ", length: " + Math.trunc(length));
}

document.addEventListener("scroll", onScroll)

window.onload = onScroll