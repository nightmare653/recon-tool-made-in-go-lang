document.getElementById("runTask").addEventListener("click", function () {
  const task = document.getElementById("task").value;
  const domain = document.getElementById("domain").value;
  const output = document.getElementById("output");

  if (!task || !domain) {
    output.textContent = "Please select a task and enter a domain.";
    return;
  }

  // Connect to WebSocket server
  const ws = new WebSocket("ws://localhost:8080/run?task=" + task + "&domain=" + domain);

  // On connection open
  ws.onopen = () => {
    output.textContent = "Running task: " + task + " on domain: " + domain + "\n";
  };

  // On receiving messages
  ws.onmessage = (event) => {
    output.textContent += event.data + "\n";
  };

  // On connection close
  ws.onclose = () => {
    output.textContent += "Connection closed.\n";
  };

  // On error
  ws.onerror = (error) => {
    output.textContent += "Error: " + error.message + "\n";
  };
});
