var selectedchat = "general";

class Event {
    constructor(type, payload){
        this.type = type;
        this.payload = payload;
    }
}
function routeEvent(event){
    if (event.type === undefined){
        alert('no type field in the event');
    }
    switch(event.type)
    {
        case "new_message" :
            console.log("new message")
            break;
        default :
        alert("unsupported message type");
        break;    
    }
}
function sendEvent(eventName, payload){
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
}

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
    sendEvent("send_message", newmessage.value)
  }
  return false;
}

window.onload = function () {
  document.getElementById("chatroom-selection").onclick = changeChatroom;
  document.getElementById("chatroom-message").onclick = sendMessage;

  if (window["WebSocket"]) {
    console.log("browser supports websockets");
    //connecting to websockets
    conn = new WebSocket("ws://" + document.location.host + "/ws"); // here ws is used to do the api connection
    conn.onmessage = function(evt){
        const eventData = JSON.parse(evt.data);
        const event = Object.assign(new Event, eventData);
        routeEvent(event);
    }
  } else {
    alert("browser does not support websocket");
  }
};
