function connect() {
  let url = "ws://localhost:4200/s";
  let socket = new WebSocket(url);
  socket.onopen = function (event) {
    socket.send("{ \"status\": \"OK\" }");
  };
  socket.onmessage = function (event) {
    console.log(event.data);
  }
}
