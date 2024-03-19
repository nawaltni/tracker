const WebSocket = require("ws");
const { v4: uuidv4 } = require("uuid");

const exampleSocket = new WebSocket("ws://localhost:8080/ws");

exampleSocket.onmessage = (event) => {
  console.log(event.data);
};

// exampleSocket.send("Hello Server!");

exampleSocket.onopen = () => {
  console.log("Connection established...");

  // send get_posts message to server as a json message
  exampleSocket.send(
    JSON.stringify({
      type: "get_positions",
      correlation_id: generateUUID(),
      data: {
        user_id: "123",
      },
    })
  );

  console.log("Message sent...");

  //wait for 2 seconds and send another message
  setTimeout(() => {
    exampleSocket.send(
      JSON.stringify({
        type: "stream_positions",
        correlation_id: generateUUID(),
        data: {
          user_id: "123",
        },
      })
    );
  }, 2000);
};

exampleSocket.onclose = () => {
  console.log("Connection closed...");
};

exampleSocket.onerror = (error) => {
  console.log(`Error: ${error}`);
};

const generateUUID = () => {
  return uuidv4();
};
