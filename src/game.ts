var data: any;
var tiles: HTMLElement[];
var conn: WebSocket;

function getGameState() {
  let req = new XMLHttpRequest();

  req.addEventListener("load", function(event) {
    data = JSON.parse(req.responseText);
    console.log(data);

    conn = new WebSocket("ws://" + document.location.host + `/websocket/${data.roomCode}`)
    console.log("ws://" + document.location.host + "/websocket");
    conn.onopen = function (evt) {
      console.log("Connection established.");
    };
    conn.onclose = function (evt) {
      console.log("Connection closed.");
    };
    conn.onmessage = updateState;
    
    updateView(data);
  });

  console.log(window.location.pathname);
  req.open("GET", "/api/v1/rooms" + window.location.pathname);
  req.send();
}

function updateView(data: any[]) {
  // Get the templates from the DOM
  let roomCodeTemplate = $("#room-code-template").html();
  let headerTemplate = $("#game-state-header-template").html();
  let boardTemplate = $("#game-state-board-template").html();

  // Create render functions for the templates with doT.template
  let roomCodeRenderFunction = doT.template(roomCodeTemplate);
  let headerRenderFunction = doT.template(headerTemplate);
  let boardRenderFunction = doT.template(boardTemplate);

  // Use the render functions to render the data
  let roomCodeRendered = roomCodeRenderFunction(data);
  let headerRendered = headerRenderFunction(data);
  let boardRendered = boardRenderFunction(data);

  // Insert the rendered results into the DOM
  $("#room-code").html(roomCodeRendered);
  $("#game-state-header").html(headerRendered);
  $("#game-state-board").html(boardRendered);

  // update tiles variable
  tiles = $(".word-tile").toArray();
  tiles.forEach((tile, index) => {
    if (!tile.classList.contains("unrevealed")) {
      tile.removeEventListener("click", () => sendMove(index));
    } else {
      tile.addEventListener("click", () => sendMove(index));
    }
  });
}

function operativeView() {
  let operativeToggle = $("#operative")[0] as HTMLInputElement;
  let spymasterToggle = $("#spymaster")[0] as HTMLInputElement;
  console.log(!operativeToggle.checked);
  if (!operativeToggle.checked) {
    console.log("Switching to operative view");
    operativeToggle.checked = true;
    spymasterToggle.checked = false;
    tiles.forEach((tile, index) => {
      if (tile.classList.contains(`${data.words[index].identity}-unrevealed`)) {
        tile.classList.add("unrevealed");
        tile.classList.remove(`${data.words[index].identity}-unrevealed`);
        tile.addEventListener("click", () => sendMove(index));
      }
    });
    $("#end-turn-btn").show();
  };
}
function spymasterView() {
  let operativeToggle = $("#operative")[0] as HTMLInputElement;
  let spymasterToggle = $("#spymaster")[0] as HTMLInputElement;
  console.log(!spymasterToggle.checked);
  if (!spymasterToggle.checked) {
    console.log("Switching to spymaster view");
    spymasterToggle.checked = true;
    operativeToggle.checked = false;
    tiles.forEach((tile, index) => {
      if (tile.classList.contains("unrevealed")) {
        tile.classList.add(`${data.words[index].identity}-unrevealed`);
        tile.classList.remove("unrevealed");
        tile.removeEventListener("click", () => sendMove(index));
      }
    });
    $("#end-turn-btn").hide();
  };
}

function sendMove(index: number) {
  conn.send(JSON.stringify(data.words[index]));
  console.log("SENDING " + JSON.stringify(data.words[index]))
}

function updateState(event: MessageEvent) {
  if (event.data != null) {
    console.log("EVENT " + JSON.stringify(event));
    console.log("EVENT DATA " + event.data);
    try {
      var dataParsed = JSON.parse(event.data);
      console.log(dataParsed);
      data = dataParsed;
      updateView(data);
    } catch(e) {
      console.log(e);
    }
  }
}

function copyRoomLinkToClipboard() {
  var $temp = $("<input>");
  $("body").append($temp);
  $temp.val(location.host + location.pathname).select();
  document.execCommand("copy");
  $temp.remove();
  
  (<any>$("#copy-btn")).focus().popover({
    trigger: 'focus',
  })
}

$(function () {
  (<any>$('[data-toggle="popover"]')).popover()
})

$("#operative").parent().on("click", operativeView);
$("#spymaster").parent().on("click", spymasterView);
$("#copy-btn").on("click", copyRoomLinkToClipboard);

getGameState();