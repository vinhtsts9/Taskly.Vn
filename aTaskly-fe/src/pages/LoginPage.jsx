import React, { useState, useContext, useEffect } from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { AuthContext } from "../context/AuthContext";
import "./AuthForm.css";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { login } = useContext(AuthContext);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Hiển thị thông báo nếu được chuyển hướng từ trang đăng ký
    if (location.state?.message) {
      setSuccessMessage(location.state.message);
      // Xóa state khỏi location để không hiển thị lại khi refresh
      window.history.replaceState({}, document.title);
    }
  }, [location.state]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      // Hàm login từ context sẽ gọi API và xử lý state
      await login({ email: email.trim().toLowerCase(), password: password });
      navigate("/"); // Điều hướng về trang chủ sau khi đăng nhập thành công
    } catch (err) {
      // Hàm login sẽ ném lỗi nếu API thất bại
      setError(err.message || "Login failed. Please check your credentials.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-form-container">
        <h1>Đăng nhập vào Taskly</h1>
        {successMessage && <p className="success-message">{successMessage}</p>}
        <form onSubmit={handleSubmit}>
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={loading}
          />

          <label htmlFor="password">Password</label>
          <div className="password-input-container">
            <input
              type={showPassword ? "text" : "password"}
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={loading}
            />
            <span
              className="password-toggle-icon"
              onClick={() => setShowPassword(!showPassword)}
            >
              {showPassword ? (
                <i className="fas fa-eye-slash"></i>
              ) : (
                <i className="fas fa-eye"></i>
              )}
            </span>
          </div>

          {error && <p className="error-message">{error}</p>}

          <button type="submit" className="auth-button" disabled={loading}>
            {loading ? "Đang đăng nhập..." : "Đăng nhập"}
          </button>
        </form>
        <div className="form-footer">
          <span>Chưa phải là thành viên? </span>
          <Link to="/register">Tham gia ngay</Link>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
