"use strict";
function connect() {
    var url = "ws://localhost:4200/s";
    var socket = new WebSocket(url);
    socket.onopen = function (event) {
        socket.send("{ \"status\": \"OK\" }");
    };
    socket.onmessage = function (event) {
        console.log(event.data);
    };
}
