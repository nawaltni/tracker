const WebSocket = require("ws");

const exampleSocket = new WebSocket("ws://localhost:8080/ws");

exampleSocket.onmessage = (event) => {
  console.log(event.data);
};

// exampleSocket.send("Hello Server!");

exampleSocket.onopen = () => {
  console.log("Connection established...");
};

exampleSocket.onclose = () => {
  console.log("Connection closed...");
};

exampleSocket.onerror = (error) => {
  console.log(`Error: ${error}`);
};
