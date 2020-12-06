function getGameState() {
  let req = new XMLHttpRequest();

  req.addEventListener("load", function(event) {
    let data = JSON.parse(req.responseText);
    console.log(data);
    
    updateView(data);
  });

  console.log(window.location.pathname);
  req.open("GET", "/api/v1/rooms" + window.location.pathname);
  req.send();
}

function updateView(data: any[]) {
  // Get the templates from the DOM
  let headerTemplate = $("#game-state-header-template").html();
  let boardTemplate = $("#game-state-board-template").html();

  // Create render functions for the templates with doT.template
  let headerRenderFunction = doT.template(headerTemplate);
  let boardRenderFunction = doT.template(boardTemplate);

  // Use the render functions to render the data
  let headerRendered = headerRenderFunction(data);
  let boardRendered = boardRenderFunction(data);

  // Insert the rendered results into the DOM
  $("#game-state-header").html(headerRendered);
  $("#game-state-board").html(boardRendered);
}

$(function () {
  (<any>$('[data-toggle="popover"]')).popover()
  })

  getGameState();