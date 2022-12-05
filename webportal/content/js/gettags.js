

function getTags(setTags) {
    fetch('/tags', {
        method: 'GET',
        credentials: 'include'
    })
    .then((response) => response.json())
    .then((data) => {
        console.log(data)
        setTags(data)
    });
}