import React, { useEffect } from "react";

const AGENT_ID = import.meta.env.AGENT_ID || "";
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

  return (
    <div style={{ position: "fixed", bottom: 32, right: 32, zIndex: 1000 }}>
      <elevenlabs-convai agent-id={AGENT_ID}></elevenlabs-convai>
    </div>
  );
};

export default AIChatbot;
