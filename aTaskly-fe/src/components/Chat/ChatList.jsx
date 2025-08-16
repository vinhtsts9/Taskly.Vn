import React, { useState, useEffect, useContext } from 'react';
import './ChatList.css';
import { AuthContext } from '../../context/AuthContext';
import { apiGetAuth } from '../../utils/api';

const ChatList = ({ onSelectThread, initialRoomId }) => {
  const [threads, setThreads] = useState([]);
  const [loading, setLoading] = useState(true);
  const { currentUser } = useContext(AuthContext);

  useEffect(() => {
    const fetchThreads = async () => {
      if (!currentUser) return;

      setLoading(true);
      try {
        const response = await apiGetAuth('/chat/rooms');
        let roomsData = response.data ?? response;

        const transformedThreads = (roomsData || []).map(room => {
          const isCurrentUserUser1 = room.user1_id === currentUser.id;
          const otherUser = isCurrentUserUser1
            ? { id: room.user2_id, name: room.user2_name, avatar: room.user2_profile_pic }
            : { id: room.user1_id, name: room.user1_name, avatar: room.user1_profile_pic };

          return {
            id: room.id,
            roomId: room.id,
            type: 'chat',
            otherUser,
            lastMessage: {
              content: room.last_message || '',
              time: room.last_message_time || room.created_at
            }
          };
        });

        setThreads(transformedThreads);
      } catch (error) {
        console.error('Error fetching chat threads:', error);
        setThreads([]);
      } finally {
        setLoading(false);
      }
    };

    fetchThreads();
  }, [currentUser]);

  return (
    <div className="chat-list">
      <div className="chat-list-header">
        <h3>Tin nhắn</h3>
      </div>
      <div className="chat-threads">
        {loading ? (
          <p>Đang tải danh sách chat...</p>
        ) : threads.length === 0 ? (
          <p>Chưa có cuộc trò chuyện nào</p>
        ) : (
          threads.map(thread => (
            <div
              key={thread.id}
              className="chat-thread-item"
              onClick={() => onSelectThread(thread)}
            >
              <div className="thread-avatar">
                {thread.otherUser.avatar
                  ? <img src={thread.otherUser.avatar} alt="avatar" />
                  : <div className="avatar-placeholder">{thread.otherUser.name.charAt(0)}</div>
                }
              </div>
              <div className="thread-content">
                <div className="thread-name">{thread.otherUser.name}</div>
                <div className="thread-last-message">{thread.lastMessage.content || 'Không có tin nhắn'}</div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ChatList;
