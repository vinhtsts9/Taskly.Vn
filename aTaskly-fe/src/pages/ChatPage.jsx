import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import './ChatPage.css';
import ChatContainer from '../components/Chat/ChatContainer';
import ErrorBoundary from '../components/ErrorBoundary';
import websocketService from '../services/websocketService';

const ChatPage = () => {
  const [searchParams] = useSearchParams();
  const [initialRoomId, setInitialRoomId] = useState(null);

  useEffect(() => {
    // Lấy room_id từ URL parameters
    const roomId = searchParams.get('room');
    const tempRoom = searchParams.get('temp_room');
    
    if (roomId) {
      // Lấy thông tin seller từ localStorage
      const sellerInfoStr = localStorage.getItem('temp_seller_info');
      let sellerInfo = null;
      
      if (sellerInfoStr) {
        try {
          sellerInfo = JSON.parse(sellerInfoStr);
          // Xóa thông tin tạm thời sau khi đã sử dụng
          localStorage.removeItem('temp_seller_info');
        } catch (error) {
          console.error('Error parsing seller info:', error);
        }
      }
      
      setInitialRoomId({
        roomId,
        sellerInfo
      });
    } else if (tempRoom === 'true') {
      // Trường hợp temp_room - chỉ lấy seller info, không có room_id
      const sellerInfoStr = localStorage.getItem('temp_seller_info');
      let sellerInfo = null;
      
      if (sellerInfoStr) {
        try {
          sellerInfo = JSON.parse(sellerInfoStr);
          // KHÔNG xóa localStorage vì cần dùng để tạo room sau
        } catch (error) {
          console.error('Error parsing seller info:', error);
        }
      }
      
      setInitialRoomId({
        roomId: null, // Không có room_id
        sellerInfo,
        isTempRoom: true
      });
    }

    // Connect to WebSocket when component mounts
    //websocketService.connect();

    // Disconnect from WebSocket when component unmounts
    // return () => {
    //   websocketService.disconnect();
    // };
  }, [searchParams]);

  return (
    <ErrorBoundary>
      <div className="chat-page">
        <ChatContainer initialRoomId={initialRoomId} />
      </div>
    </ErrorBoundary>
  );
};

export default ChatPage; 