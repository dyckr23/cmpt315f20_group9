"use strict";
$("#room-code-input").on("keydown", function search(e) {
    if (e.key === 'Enter') {
        var value = $(this).val();
        let req = new XMLHttpRequest();
        req.addEventListener("load", function (event) {
            let data = JSON.parse(req.responseText);
            console.log(data);
            window.location.href = "/" + value;
        });
        req.open("GET", "/api/v1/rooms/" + value);
        req.send();
        //console.log(window.location.pathname);
    }
});
