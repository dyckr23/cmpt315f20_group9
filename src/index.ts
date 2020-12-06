$("#room-code-input").on("keydown", function search(e) {
    if (e.key === 'Enter') {
        var value = <string>$(this).val();

        let req = new XMLHttpRequest();

        req.addEventListener("load", function(event) {
          let data = JSON.parse(req.responseText);
          console.log(data);
        });
      
        
        req.open("GET", "/api/v1/rooms/" + value);
        req.send();

        console.log(window.location.pathname);

        window.location.href = "/" + value;
    }
});