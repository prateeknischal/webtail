var open_connection = function(file) {
  console.log(file)
  $("#filename").html("File: " + file)
  var ws;
  if (window.WebSocket === undefined) {
      $("#container").append("Your browser does not support WebSockets");
      return;
  } else {
      ws = initWS(file);
  }
}
function initWS(file) {
  var socket = new WebSocket("ws://"+window.location.hostname+":" + window.location.port + "/ws/" + file),
  container = $("#container")
  container.html("")
  socket.onopen = function() {
      container.append("<p><b>Tailing file: " + file + "</b></p>");
  };
  socket.onmessage = function (e) {
      container.append("<p>" + e.data.trim() + "</p>");
  }
  socket.onclose = function () {
      container.append("<p>Connection Closed to WebSocket, tail stopped</p>");
  }
  return socket;
}
