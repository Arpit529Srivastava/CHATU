var selectedchat = "general";

function changeChatroom() {
  var newchat = document.getElementById("chatroom");
  if (newchat != null && newchat.value != selectedchat) {
    console.log(newchat);
  }
  return false; // prevent form navigation
}

function sendMessage() {
  var newmessage = document.getElementById("message");
  if (newmessage != null) {
    console.log(newmessage);
  }
  return false;
}

window.onload = function () {
  document.getElementById("chatroom-selection").onclick = changeChatroom;
  document.getElementById("chatroom-message").onclick = sendMessage;

  if (window["WebSocket"]) {
    console.log("browser supports websockets");
  } else {
    alert("browser does not support websocket");
  }
};
