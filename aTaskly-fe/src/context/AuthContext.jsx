import React, { createContext, useState, useEffect, useCallback } from 'react';
import { apiPost, apiGetAuth, apiPostAuth } from '../utils/api'; 
import websocketService from '../services/websocketService';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [loading, setLoading] = useState(true); // Thêm state loading để chờ kiểm tra auth

  const checkAuthStatus = useCallback(async () => {
    setLoading(true);
    try {
      // Endpoint này sẽ đọc httpOnly cookie và trả về user nếu hợp lệ
      const user = await apiGetAuth('/users/me'); 
      try { console.log('[AuthContext] /users/me ->', user); } catch(_) {}
      setCurrentUser(user);
    } catch (error) {
      // Nếu có lỗi (ví dụ: 401 Unauthorized), nghĩa là chưa đăng nhập
      try { console.warn('[AuthContext] /users/me error ->', error); } catch(_) {}
      setCurrentUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    checkAuthStatus();

    // Lắng nghe sự kiện đăng xuất toàn cục
    const handleAuthFailure = () => {
      console.log("Auth failure detected, logging out.");
      setCurrentUser(null); 
    };

    window.addEventListener('auth-failure', handleAuthFailure);

    // Dọn dẹp listener khi component bị unmount
    return () => {
      window.removeEventListener('auth-failure', handleAuthFailure);
    };
  }, [checkAuthStatus]);

  const login = async (credentials) => {
    // Backend sẽ set httpOnly cookie, frontend chỉ cần lấy thông tin user
    // Thay đổi apiPost -> apiPostAuth để trình duyệt chấp nhận Set-Cookie header
    const user = await apiPostAuth('/users/login', credentials);
    try { console.log('[AuthContext] login ->', user); } catch(_) {}
    setCurrentUser(user);
    await checkAuthStatus();   // ✅ confirm lại bằng cookie
    return user;
  };

  const logout = async () => {
    try {
      // Đóng WebSocket trước khi logout
      websocketService.disconnect();
      
      // Gọi API để backend xóa httpOnly cookie
      await apiPostAuth('/users/logout', {});
    } catch (error) {
      console.error('Error during logout:', error);
    } finally {
      try { console.log('[AuthContext] logout -> clearing user'); } catch(_) {}
      setCurrentUser(null);
      await checkAuthStatus();   // ✅ confirm lại bằng cookie
    }
  };

  // Thêm hàm checkAuth lại để các component khác có thể gọi khi cần
  const refreshUser = useCallback(() => {
    checkAuthStatus();
  }, [checkAuthStatus]);

  // Chỉ render children khi đã kiểm tra xong auth
  if (loading) {
    return <div>Đang tải ứng dụng...</div>; // Hoặc một spinner đẹp hơn
  }

  return (
    <AuthContext.Provider value={{ currentUser, login, logout, refreshUser, isAuthenticated: !!currentUser }}>
      {children}
    </AuthContext.Provider>
  );
}; 