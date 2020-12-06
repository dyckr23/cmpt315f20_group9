$("#room-code-input").on("keydown", function search(e) {
    if (e.key === 'Enter') {
        var value = <string>$(this).val();
        window.location.href = "rooms/" + value;
    }
});