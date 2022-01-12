const {WebSocket} = require('ws')
let ws;
let setting_room = "10";
const ws_url = "ws://localhost:4000/connect/";
function ws_connect(){
    ws = new WebSocket(ws_url, setting_room, {
        perMessageDeflate : true
    });

  ws.on('open', function open() {
          console.log("OPEN!!")
  }); 
  ws.on('message', function incoming(message) {
        console.log("received: %s", message)
        ws.send(message)
  }); 
  ws.on('ping', function get_ping(){
   console.log("Get PING")
   //ws.pong()
  });
  ws.on('close', function close() {
    console.log("close")
  }); 
  ws.on('error', function errorhandler(err){
      console.log("error")
  })
}
ws_connect()