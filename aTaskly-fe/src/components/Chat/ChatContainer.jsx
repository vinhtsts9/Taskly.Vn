import React, { useState, useEffect } from "react";
import "./ChatContainer.css";
import ChatList from "./ChatList";
import ChatWindow from "./ChatWindow";

const ChatContainer = ({ initialRoomId }) => {
  const [selectedThread, setSelectedThread] = useState(null);

  const handleSelectThread = (thread) => {
    setSelectedThread(thread);
  };

  // Khi initialRoomId thay đổi, tạo thread tạm hoặc thread thật
  useEffect(() => {
    if (!initialRoomId || selectedThread) return;

    if (initialRoomId.isTempRoom && initialRoomId.sellerInfo) {
      // Thread tạm thời
      const tempThread = {
        roomId: null,
        type: "temp_chat",
        isTempRoom: true,
        otherUser: initialRoomId.sellerInfo,
      };
      setSelectedThread(tempThread);
    } else if (initialRoomId.roomId && initialRoomId.sellerInfo) {
      // Thread thật
      const realThread = {
        roomId: initialRoomId.roomId,
        type: "chat",
        otherUser: initialRoomId.sellerInfo,
      };
      setSelectedThread(realThread);
    }
  }, [initialRoomId, selectedThread]);

  return (
    <div className="chat-container">
      <ChatList
        onSelectThread={handleSelectThread}
        initialRoomId={initialRoomId}
      />
      <ChatWindow
        selectedThread={selectedThread}
        setSelectedThread={setSelectedThread}
      />
    </div>
  );
};

export default ChatContainer;
