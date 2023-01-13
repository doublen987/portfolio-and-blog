
function init3dViewer(tagName, options = {}) {
    // var scene = new THREE.Scene();
    // var camera = new THREE.PerspectiveCamera(75, domElement.innerWidth/domElement.innerHeight, 0.1, 1000);

    // var renderer = new THREE.WebGLRenderer();
    // renderer.setSize(domElement.innerWidth, domElement.innerHeight);
    // domElement.appendChild(renderer.domElement)

    // var geometry = new THREE.BoxGemoetry(1,1,1);
    // var material = new THREE.MeshBasicMaterial({color: 0x00ff00});

    // camera.position.z = 5;

    // var animate = function() {
    //     requestAnimationFrame(animate);
    //     renderer.render(scene,animate)
    // }

    // animate();
    var stl_viewer=new StlViewer(document.getElementById(tagName));
    return stl_viewer

}