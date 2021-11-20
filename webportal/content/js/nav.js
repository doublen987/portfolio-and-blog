let modalhidden = true;
function initMobileMenu() {
    
    var icon = document.getElementsByClassName("nav-icon-container")[0];
    icon.addEventListener("click", () => {
        var modal = document.getElementById("modal");
        var iconholder = document.getElementsByClassName("nav-icon-container")[0];
        if(!modalhidden) {
            modal.classList.remove("modal-visible");
            modal.classList.add("modal-hidden");
            iconholder.classList.remove("nav-icon-container-clicked");
            iconholder.classList.add("nav-icon-container-not-clicked");
            modalhidden = true;
        } else {
            modal.classList.remove("modal-hidden");
            modal.classList.add("modal-visible");
            iconholder.classList.remove("nav-icon-container-not-clicked");
            iconholder.classList.add("nav-icon-container-clicked");
            
            modalhidden = false;
        }
    });
}