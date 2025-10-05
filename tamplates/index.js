const ws = new WebSocket("ws://localhost:8080/ws");
const chatBox = document.getElementById("chat-box");

ws.onmessage = (event) => {
  const msg = document.createElement("div");
  msg.textContent = event.data;

  // classify as sent or received
  if (event.data.startsWith("🧑")) {
    msg.classList.add("message", "sent");
  } else {
    msg.classList.add("message", "received");
  }

  chatBox.appendChild(msg);
  chatBox.scrollTop = chatBox.scrollHeight;
};

function sendMessage() {
  const input = document.getElementById("msg");
  if (input.value.trim() === "") return;

  ws.send(input.value);
  input.value = "";
}

document.getElementById("msg").addEventListener("keydown", (e) => {
  if (e.key === "Enter") sendMessage();
});
