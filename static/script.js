function setActiveHeaderItem(id) {
    var items = document.getElementsByName("header");
    items.forEach(item => {
        item.classList.remove("active")
    });

    var item = document.getElementById(id);
    item.classList.add("active")
}

function toggleSpoiler(id) {
    var spoilers = document.getElementsByName("spoiler");
    spoilers.forEach(spoiler => {
        var baseId = spoiler.id.substring("spoiler_".length);
        var button = document.getElementById("button_" + baseId);
        if (baseId != id && spoiler.style.display === "block") {
            spoiler.style.display = "none";
            button.innerText = "Show";
        }
    });

    var spoiler = document.getElementById("spoiler_" + id);
    var button = document.getElementById("button_" + id);
    if (spoiler.style.display === "block") {
        spoiler.style.display = "none";
        button.innerText = "Show";
    } else {
        spoiler.style.display = "block";
        button.innerText = "Hide";
    }
};

function showContext() {
    var context = document.getElementById("context");
    if (document.getElementById("describe").checked == true) {
        context.style.display = "block";
    } else {
        context.style.display = "none";
    }
}

function selectFeed() {
    option = document.cookie
        .split('; ')
        .find(row => row.startsWith('feed='))
        ?.split('=')[1];
    if (option) {
        document.getElementById("feed").value = option;
    }
}

function changeFeed() {
    var option = document.getElementById("feed").value;
    document.cookie = `feed=${option}; path=/; SameSite=Strict`;
    // reload without POST DATA
    window.location = window.location.href;
}
