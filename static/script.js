function setActiveHeaderItem(id) {
    var items = document.getElementsByName("header");
    items.forEach(item => {
        item.classList.remove("active")
    });
    var item = document.getElementById(id);
    item.classList.add("active")
}

function deleteItem(id) {
    document.getElementById("item_" + id).style.display = "none";
    document.getElementById("hr_" + id).style.display = "none";
    fetch("?delete=" + id)
        .then(data => console.log('Success:', data))
        .catch(error => console.error('Error:', error));
};

function toggleSpoiler(id) {
    var spoilers = document.getElementsByName("spoiler");
    spoilers.forEach(spoiler => {
        var baseId = spoiler.id.substring("spoiler_".length);
        var button = document.getElementById("spoilerButton_" + baseId);
        if (baseId != id && spoiler.style.display === "block") {
            spoiler.style.display = "none";
            button.innerText = "Show";
        }
    });
    var spoiler = document.getElementById("spoiler_" + id);
    var button = document.getElementById("spoilerButton_" + id);
    if (spoiler.style.display === "block") {
        spoiler.style.display = "none";
        button.innerText = "Show";
    } else {
        spoiler.style.display = "block";
        button.innerText = "Hide";
    }
};

function toggleContext() {
    var text = document.getElementById("text");
    if (document.getElementById("context").checked == true) {
        text.style.display = "block";
    } else {
        text.style.display = "none";
    }
}

function selectFeed() {
    var option = document.getElementById("header-feed-selector").value;
    document.cookie = `feed=${option}; path=/; SameSite=Strict`;
    // reload without POST DATA
    window.location = window.location.href;
}
