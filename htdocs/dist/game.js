"use strict";
var data;
var tiles;
var callbacks;
var conn;
function getGameState() {
    let req = new XMLHttpRequest();
    req.addEventListener("load", function (event) {
        data = JSON.parse(req.responseText);
        console.log(data);
        conn = new WebSocket("ws://" + document.location.host + `/websocket/${data.roomCode}`);
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
function updateView(data) {
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
    // Check game status
    // If ongoing, everything is normal and clickable
    if (!callbacks) {
        callbacks = Array(25).fill(null);
    }
    if (data.status == "ongoing") {
        tiles.forEach((tile, index) => {
            callbacks[index] = (callbacks[index] != null) ? callbacks[index] : () => sendMove(index);
            if (!tile.classList.contains("unrevealed")) {
                tile.removeEventListener("click", callbacks[index]);
            }
            else {
                tile.addEventListener("click", callbacks[index]);
            }
        });
    }
    // If not, disable tile clicking and "end turn" button
    else {
        tiles.forEach((tile, index) => {
            tile.removeEventListener("click", callbacks[index]);
            tile.classList.add("disabled");
        });
        $("#end-turn-btn").prop("disabled", true);
    }
    console.log(callbacks);
}
function operativeView() {
    let operativeToggle = $("#operative")[0];
    let spymasterToggle = $("#spymaster")[0];
    console.log(!operativeToggle.checked);
    if (!operativeToggle.checked) {
        console.log("Switching to operative view");
        operativeToggle.checked = true;
        spymasterToggle.checked = false;
        tiles.forEach((tile, index) => {
            if (tile.classList.contains(`${data.words[index].identity}-unrevealed`)) {
                tile.addEventListener("click", callbacks[index]);
                tile.classList.add("unrevealed");
                tile.classList.remove(`${data.words[index].identity}-unrevealed`);
            }
        });
        $("#end-turn-btn").show();
    }
    ;
}
function spymasterView() {
    let operativeToggle = $("#operative")[0];
    let spymasterToggle = $("#spymaster")[0];
    console.log(!spymasterToggle.checked);
    if (!spymasterToggle.checked) {
        console.log("Switching to spymaster view");
        spymasterToggle.checked = true;
        operativeToggle.checked = false;
        tiles.forEach((tile, index) => {
            if (tile.classList.contains("unrevealed")) {
                tile.removeEventListener("click", callbacks[index]);
                tile.classList.add(`${data.words[index].identity}-unrevealed`);
                tile.classList.remove("unrevealed");
            }
        });
        $("#end-turn-btn").hide();
    }
    ;
}
function sendMove(index) {
    conn.send(JSON.stringify(data.words[index]));
    console.log("SENDING " + JSON.stringify(data.words[index]));
}
function updateState(event) {
    if (event.data != null) {
        console.log("EVENT " + JSON.stringify(event));
        console.log("EVENT DATA " + event.data);
        try {
            var dataParsed = JSON.parse(event.data);
            console.log(dataParsed);
            data = dataParsed;
            updateView(data);
        }
        catch (e) {
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
    $("#copy-btn").focus().popover({
        trigger: 'focus',
    });
}
$(function () {
    $('[data-toggle="popover"]').popover();
});
$("#operative").parent().on("click", operativeView);
$("#spymaster").parent().on("click", spymasterView);
$("#copy-btn").on("click", copyRoomLinkToClipboard);
getGameState();
