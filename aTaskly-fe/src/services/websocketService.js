import { apiGet, apiPost } from "../utils/api"; // Giả sử api.js có export API_BASE_URL
const wsUrl = import.meta.env.VITE_WS_URL || "";
let socket = null;
const messageListeners = new Set();
let pingInterval = null;

const connect = () => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("WebSocket is already connected.");
    return;
  }

  console.log("Connecting to WebSocket:", wsUrl);
  socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    console.log("WebSocket connected successfully.");
    // Gửi ping định kỳ mỗi 50s
    if (pingInterval) clearInterval(pingInterval);
    pingInterval = setInterval(() => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ action: "ping" }));
      }
    }, 50000); // 50s
  };

  socket.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      for (const listener of messageListeners) {
        listener(message);
      }
    } catch (error) {
      console.error("Error parsing incoming message:", error);
    }
  };

  socket.onclose = () => {
    console.log("WebSocket disconnected.");
    if (pingInterval) clearInterval(pingInterval);
    socket = null;
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };
};

const disconnect = () => {
  if (socket) {
    socket.close();
  }
};

const send = (message) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message));
  } else {
    console.log("WebSocket not ready, attempting to reconnect...");
    connect();
    // Wait a bit for connection to establish
    setTimeout(() => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
      } else {
        console.error(
          "WebSocket is not connected or ready after reconnect attempt."
        );
      }
    }, 1000);
  }
};

const joinRoom = (roomId) => {
  send({
    action: "join",
    room_id: roomId,
  });
};

const leaveRoom = (roomId) => {
  send({
    action: "leave",
    room_id: roomId,
  });
};

const sendChatMessage = (roomId, receiverId, content) => {
  send({
    action: "send_message",
    room_id: roomId,
    receiver_id: receiverId,
    content: content,
  });
};

const addMessageListener = (callback) => {
  messageListeners.add(callback);
};

const removeMessageListener = (callback) => {
  messageListeners.delete(callback);
};

const isConnected = () => {
  return socket && socket.readyState === WebSocket.OPEN;
};

const websocketService = {
  connect,
  disconnect,
  joinRoom,
  leaveRoom,
  sendChatMessage,
  addMessageListener,
  removeMessageListener,
  isConnected,
};

export default websocketService;
