import React, { useEffect } from "react";

const AGENT_ID = import.meta.env.VITE_AGENT_ID || "";
const AIChatbot = () => {
  useEffect(() => {
    // Chỉ thêm script nếu chưa có
    if (!document.getElementById("elevenlabs-convai-script")) {
      const script = document.createElement("script");
      script.src = "https://elevenlabs.io/convai-widget/index.js";
      script.async = true;
      script.type = "text/javascript";
      script.id = "elevenlabs-convai-script";
      document.body.appendChild(script);
    }
  }, []);

  if (!AGENT_ID) {
    return (
      <div
        style={{
          position: "fixed",
          bottom: 32,
          right: 32,
          zIndex: 1000,
          background: "#fff3cd",
          color: "#856404",
          padding: "12px 20px",
          borderRadius: 8,
          boxShadow: "0 2px 8px rgba(0,0,0,0.12)",
        }}
      >
        <b>Chưa cấu hình VITE_AGENT_ID!</b>
        <br />
        Vui lòng thêm VITE_AGENT_ID vào file .env để chatbot hoạt động.
      </div>
    );
  }
  return (
    <div style={{ position: "fixed", bottom: 32, right: 32, zIndex: 1000 }}>
      <elevenlabs-convai agent-id={AGENT_ID}></elevenlabs-convai>
    </div>
  );
};

export default AIChatbot;
