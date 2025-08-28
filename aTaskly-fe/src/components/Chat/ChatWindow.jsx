import React, {
  useState,
  useEffect,
  useRef,
  useContext,
  useCallback,
  useLayoutEffect,
} from "react";
import "./ChatWindow.css";
import { AuthContext } from "../../context/AuthContext";
import websocketService from "../../services/websocketService";
import { apiGetAuth, apiPostAuth } from "../../utils/api";

const ChatWindow = ({ selectedThread, setSelectedThread }) => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [chatPartner, setChatPartner] = useState(null);
  const [historyCursor, setHistoryCursor] = useState(null);
  const [hasMoreHistory, setHasMoreHistory] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [isSending, setIsSending] = useState(false); // State để ngăn gửi tin nhắn trùng lặp

  const messagesContainerRef = useRef(null);
  const prevScrollHeightRef = useRef(null);
  const messagesEndRef = useRef(null);
  const { currentUser } = useContext(AuthContext);
  const currentRoomId = useRef(null);

  // Logic to load more messages
  const handleLoadMore = useCallback(async () => {
    if (
      isLoadingMore ||
      !hasMoreHistory ||
      !selectedThread ||
      selectedThread.isTempRoom
    )
      return;
    setIsLoadingMore(true);
    try {
      let url = `/chat/history/${selectedThread.roomId}`;
      if (historyCursor) url += `?cursor=${historyCursor}`;
      const res = await apiGetAuth(url);
      // Backend trả về DESC, phải reverse lại
      const olderMessages = (res || []).reverse();
      if (olderMessages.length > 0) {
        if (messagesContainerRef.current) {
          prevScrollHeightRef.current =
            messagesContainerRef.current.scrollHeight;
        }
        // Loại bỏ trùng lặp
        setMessages((prev) => {
          const existingIds = new Set(prev.map((m) => m.id));
          const uniqueOlder = olderMessages.filter(
            (m) => !existingIds.has(m.id)
          );
          return [...uniqueOlder, ...prev];
        });
        // Cursor là sent_at của phần tử CUỐI CÙNG (cũ nhất)
        setHistoryCursor(olderMessages[olderMessages.length - 1].sent_at);
      } else {
        setHasMoreHistory(false);
      }
    } catch (error) {
      console.error("Failed to load older messages:", error);
    } finally {
      setIsLoadingMore(false);
    }
  }, [isLoadingMore, hasMoreHistory, historyCursor, selectedThread]);
  // Scroll event handler
  const handleScroll = useCallback(() => {
    const container = messagesContainerRef.current;
    if (!container) return;
    if (
      container.scrollTop <= 10 &&
      !isLoadingMore &&
      hasMoreHistory &&
      selectedThread &&
      !selectedThread.isTempRoom
    ) {
      handleLoadMore();
    }
  }, [handleLoadMore, isLoadingMore, hasMoreHistory, selectedThread]);

  // Effect to attach/detach scroll listener
  useEffect(() => {
    const container = messagesContainerRef.current;
    if (container) {
      container.addEventListener("scroll", handleScroll);
      return () => container.removeEventListener("scroll", handleScroll);
    }
  }, [handleScroll]);

  // Effect to reset state when changing threads
  useEffect(() => {
    if (currentRoomId.current) {
      websocketService.leaveRoom(currentRoomId.current);
    }
    setMessages([]);
    setChatPartner(null);
    setHistoryCursor(null);
    setHasMoreHistory(true);
    prevScrollHeightRef.current = null;

    if (selectedThread) {
      // Nếu là temp room, không cần join WebSocket room
      if (selectedThread.isTempRoom) {
        currentRoomId.current = null;
        // Không gọi API lấy lịch sử vì chưa có room
        setChatPartner(selectedThread.otherUser);
        return;
      }

      // Room thật - join WebSocket và lấy lịch sử
      currentRoomId.current = selectedThread.roomId;

      // Ensure WebSocket is connected before joining room
      if (!websocketService.isConnected()) {
        websocketService.connect();
      }

      websocketService.joinRoom(selectedThread.roomId);

      const fetchInitialData = async () => {
        setIsLoadingMore(true);
        try {
          // Use user info from selectedThread instead of calling room info API
          const historyRes = await apiGetAuth(
            `/chat/history/${selectedThread.roomId}`
          );

          const history = historyRes || [];
          setMessages(history.reverse());

          if (history.length > 0) {
            setHistoryCursor(history[0].created_at);
            messagesEndRef.current?.scrollIntoView();
          } else {
            setHasMoreHistory(false);
          }

          // Use user info from selectedThread
          setChatPartner(selectedThread.otherUser);
        } catch (error) {
          console.error("Error fetching initial chat data:", error);
          setMessages([]);
          setChatPartner(null);
          setHasMoreHistory(false);
        } finally {
          setIsLoadingMore(false);
        }
      };
      fetchInitialData();
    }

    return () => {
      // Cleanup khi component unmount hoặc selectedThread thay đổi
      if (currentRoomId.current) {
        try {
          websocketService.leaveRoom(currentRoomId.current);
        } catch (error) {
          console.error("Error leaving room during cleanup:", error);
        }
      }
    };
  }, [selectedThread, currentUser.id]);

  // WebSocket message listener
  useEffect(() => {
    const handleNewMessage = (message) => {
      try {
        if (message.room_id === currentRoomId.current) {
          // Kiểm tra xem tin nhắn đã tồn tại chưa
          const messageExists = messages.some((m) => m.id === message.id);
          if (!messageExists) {
            // Nếu tin nhắn từ người khác, thêm vào danh sách
            if (message.sender_id !== currentUser?.id) {
              setMessages((prev) => [...prev, message]);
            } else {
              // Nếu tin nhắn từ chính mình, cập nhật tin nhắn tạm thời với ID thật
              setMessages((prev) =>
                prev.map((m) => {
                  if (
                    m.id.startsWith("temp_") &&
                    m.content === message.content
                  ) {
                    return message; // Thay thế tin nhắn tạm thời bằng tin nhắn thật
                  }
                  return m;
                })
              );
            }
          }
        }
      } catch (error) {
        console.error("Error handling new message:", error);
      }
    };

    websocketService.addMessageListener(handleNewMessage);

    return () => {
      try {
        websocketService.removeMessageListener(handleNewMessage);
      } catch (error) {
        console.error("Error removing message listener:", error);
      }
    };
  }, [messages, currentUser?.id]);

  // Effect to manage scroll position
  useLayoutEffect(() => {
    const container = messagesContainerRef.current;
    if (!container) return;

    const nearBottom =
      container.scrollHeight - container.scrollTop - container.clientHeight <=
      200;
    // 50px tolerance (gần cuối thì coi như ở cuối)

    if (prevScrollHeightRef.current) {
      // Sau khi load thêm tin nhắn cũ
      container.scrollTop =
        container.scrollHeight - prevScrollHeightRef.current;
      prevScrollHeightRef.current = null;
    } else if (nearBottom) {
      // Chỉ auto-scroll nếu user đang gần cuối
      container.scrollTop = container.scrollHeight;
    }
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    // Kiểm tra isSending để ngăn chặn double-click
    if (!newMessage.trim() || !selectedThread || !chatPartner || isSending)
      return;

    setIsSending(true);
    try {
      // Nếu là temp room, tạo room thật trước
      if (selectedThread.isTempRoom) {
        const response = await apiPostAuth("/chat/create-room", {
          content: newMessage,
          user2_id: chatPartner.id,
        });

        if (response && response.id) {
          // Cập nhật thread với room_id thật
          const updatedThread = {
            ...selectedThread,
            roomId: response.id,
            type: "chat",
            isTempRoom: false,
          };

          // Cập nhật selectedThread
          setSelectedThread(updatedThread);

          // Thêm tin nhắn vào danh sách local
          const newMessageObj = {
            id: messageResponse.id || `temp_${Date.now()}`,
            room_id: response.id,
            sender_id: currentUser.id,
            receiver_id: chatPartner.id,
            content: newMessage,
            sent_at: new Date().toISOString(),
            created_at: new Date().toISOString(),
          };

          setMessages((prev) => [...prev, newMessageObj]);

          // Sau khi tạo room và gửi tin nhắn thành công, mới join WebSocket room
          setTimeout(() => {
            if (websocketService.isConnected()) {
              websocketService.joinRoom(response.id);
              currentRoomId.current = response.id;

              // Cập nhật thread để không còn là temp room nữa
              const finalThread = {
                ...updatedThread,
                roomId: response.id,
                type: "chat",
                isTempRoom: false,
              };
              setSelectedThread(finalThread);

              // Reset state để có thể load lịch sử tin nhắn
              setMessages([newMessageObj]); // Chỉ giữ tin nhắn vừa gửi
              setHistoryCursor(null);
              setHasMoreHistory(true);

              // Trigger load lịch sử tin nhắn cũ (nếu có)
              handleLoadMore();
            }
          }, 500); // Delay 500ms để đảm bảo backend đã xử lý xong
        } else {
          console.error("Failed to create chat room");
          return;
        }
      } else {
        // Room đã tồn tại, gửi tin nhắn qua WebSocket
        websocketService.sendChatMessage(
          selectedThread.roomId,
          chatPartner.id,
          newMessage
        );

        // Thêm tin nhắn vào danh sách local ngay lập tức để UI responsive
        const newMessageObj = {
          id: `temp_${Date.now()}`, // ID tạm thời, sẽ được cập nhật khi nhận từ WebSocket
          room_id: selectedThread.roomId,
          sender_id: currentUser.id,
          receiver_id: chatPartner.id,
          content: newMessage,
          sent_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
        };

        setMessages((prev) => [...prev, newMessageObj]);
      }

      setNewMessage("");
    } catch (error) {
      console.error("Error sending message:", error);
    } finally {
      setIsSending(false); // Luôn reset trạng thái sau khi gửi xong
    }
  };

  if (!selectedThread) {
    return (
      <div className="chat-window empty-state">
        <p>Chọn một cuộc trò chuyện để bắt đầu</p>
      </div>
    );
  }

  return (
    <div className="chat-window">
      <div className="chat-header">
        {chatPartner ? (
          <div className="chat-user-info">
            {chatPartner.avatar ? (
              <img
                src={chatPartner.avatar}
                alt="avatar"
                className="chat-avatar"
                onError={(e) => {
                  e.target.style.display = "none";
                  e.target.nextSibling.style.display = "block";
                }}
              />
            ) : null}
            <div
              className="chat-avatar-placeholder"
              style={{ display: chatPartner.avatar ? "none" : "block" }}
            >
              {chatPartner.name.charAt(0).toUpperCase()}
            </div>
            <div className="chat-user-details">
              <h4>{chatPartner.name}</h4>
              {selectedThread.type === "order" && (
                <span className="order-badge">
                  Order #{selectedThread.orderId}
                </span>
              )}
            </div>
          </div>
        ) : (
          <div className="chat-user-info">Loading...</div>
        )}
      </div>

      <div className="messages-container" ref={messagesContainerRef}>
        {isLoadingMore && <div className="loading-spinner"></div>}
        {messages.map((message) => (
          <div
            key={message.id}
            className={`message ${
              message.sender_id === currentUser.id ? "sent" : "received"
            }`}
          >
            <div className="message-content">
              <p>{message.content}</p>
              <span className="message-time">
                {new Date(message.created_at).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </span>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form className="message-input" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Nhập tin nhắn..."
        />
        <button
          type="submit"
          disabled={
            !newMessage.trim() || !chatPartner || isSending || !selectedThread
          }
        >
          {isSending ? "Đang gửi..." : "Gửi"}
        </button>
      </form>
    </div>
  );
};

export default ChatWindow;
