import React, { useContext, useState, useEffect, useRef } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import { apiGet } from "../utils/api";
import "./Header.css";
import { isUserAdmin } from "../utils/auth";
import { registerPush } from "../utils/push";

const Header = () => {
  const { currentUser, logout, logoutLoading } = useContext(AuthContext);

  const [searchTerm, setSearchTerm] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const navigate = useNavigate();
  const searchContainerRef = useRef(null);

  // Debouncing effect for search suggestions
  useEffect(() => {
    const handler = setTimeout(() => {
      if (searchTerm.trim()) {
        fetchSuggestions(searchTerm.trim());
      } else {
        setSuggestions([]);
        setShowSuggestions(false);
      }
    }, 300); // 300ms delay

    return () => {
      clearTimeout(handler);
    };
  }, [searchTerm]);

  // Effect to handle clicks outside the search container
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        searchContainerRef.current &&
        !searchContainerRef.current.contains(event.target)
      ) {
        setShowSuggestions(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const fetchSuggestions = async (keyword) => {
    try {
      // Sửa đổi endpoint và tham số
      const data = await apiGet(
        `/gigs/search?search_term=${encodeURIComponent(keyword)}`
      );
      // Định dạng phản hồi API đã thay đổi
      const titles = data.map((item) => item.title);
      setSuggestions(titles);
      setShowSuggestions(true);
    } catch (error) {
      console.error("Failed to fetch suggestions:", error);
      setSuggestions([]);
      setShowSuggestions(false);
    }
  };

  const handleSearch = (keyword) => {
    const finalKeyword = keyword.trim();
    if (finalKeyword) {
      setSearchTerm(finalKeyword);
      setShowSuggestions(false);
      // Sửa đổi tham số chuyển hướng
      navigate(`/gigs?search_term=${encodeURIComponent(finalKeyword)}`);
    }
  };

  const handleKeyDown = (event) => {
    if (event.key === "Enter") {
      handleSearch(searchTerm);
    }
  };

  const handleEnableNotifications = async () => {
    try {
      await registerPush();
      alert("Bạn đã bật thông báo thành công!");
    } catch (error) {
      console.error("Không thể bật thông báo:", error);
      alert("Đã xảy ra lỗi khi bật thông báo. Vui lòng thử lại.");
    }
  };

  return (
    <header>
      <nav>
        <div className="logo">
          <Link to="/">
            <span>Taskly</span>
          </Link>
        </div>
        <div className="search-bar" ref={searchContainerRef}>
          <input
            type="text"
            placeholder="Bạn đang tìm kiếm dịch vụ gì?"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onFocus={() => setShowSuggestions(true)}
            onKeyDown={handleKeyDown}
          />
          <button onClick={() => handleSearch(searchTerm)}>Tìm kiếm</button>
          {showSuggestions && suggestions.length > 0 && (
            <ul className="suggestions-list">
              {suggestions.map((suggestion, index) => (
                <li key={index} onClick={() => handleSearch(suggestion)}>
                  {suggestion}
                </li>
              ))}
            </ul>
          )}
        </div>
        <ul className="nav-links">
          {currentUser ? (
            <>
              {(() => {
                try {
                  console.log("[Header] currentUser:", currentUser);
                } catch (_) {}
                return null;
              })()}
              <li>
                <Link to="/dashboard/orders">Đơn hàng</Link>
              </li>
              <li>
                <Link to="/messages" className="messages-link">
                  Tin nhắn
                </Link>
              </li>
              {/* Hiển thị link Admin nếu user có quyền */}
              {isUserAdmin(currentUser) && (
                <li>
                  <Link to="/admin">Admin</Link>
                </li>
              )}
              <li>
                <button
                  onClick={handleEnableNotifications}
                  className="notification-button"
                >
                  Bật thông báo
                </button>
              </li>
              <li>
                <span>Chào, {currentUser.names}</span>
              </li>
              <li>
                <button
                  onClick={logout}
                  className="logout-button"
                  disabled={logoutLoading}
                >
                  {logoutLoading ? "Đang xử lý..." : "Đăng xuất"}
                </button>
              </li>
            </>
          ) : (
            <>
              {/* Nếu chưa là seller thì hiện 'Trở thành người bán', nếu là seller thì hiện 'Đăng dịch vụ' */}
              {currentUser && currentUser.role ? (
                <li>
                  <Link to="/gigs/new">Đăng dịch vụ</Link>
                </li>
              ) : (
                <li>
                  <Link to="/become-a-seller">Trở thành người bán</Link>
                </li>
              )}
              <li>
                <Link to="/login">Đăng nhập</Link>
              </li>
              <li>
                <Link to="/register">
                  <button className="join-button">Tham gia</button>
                </Link>
              </li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;
