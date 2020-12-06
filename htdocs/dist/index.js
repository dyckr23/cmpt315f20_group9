"use strict";
$("#room-code-input").on("keydown", function search(e) {
    if (e.key === 'Enter') {
        var value = $(this).val();
        window.location.href = "rooms/" + value;
    }
});
