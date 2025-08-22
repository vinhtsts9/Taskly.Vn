import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import "./AuthForm.css";
import { apiPost } from "../utils/api";

const RegisterPage = () => {
  const [step, setStep] = useState(1); // 1: Enter email/phone, 2: Verify OTP, 3: Details
  const [verifyKey, setVerifyKey] = useState("");
  const [verifyType, setVerifyType] = useState("email"); // 'email' or 'phone'
  const [otp, setOtp] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [names, setNames] = useState("");
  const [userType, setUserType] = useState("buyer"); // Mặc định là 'buyer'
  const [bio, setBio] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // State và hàm cho việc kiểm tra độ mạnh mật khẩu
  const [passwordValidation, setPasswordValidation] = useState({
    minLength: false,
    hasUpper: false,
    hasLower: false,
    hasNumber: false,
    hasSpecial: false,
  });

  const validatePassword = (value) => {
    setPassword(value);
    const minLength = value.length >= 8;
    const hasUpper = /[A-Z]/.test(value);
    const hasLower = /[a-z]/.test(value);
    const hasNumber = /[0-9]/.test(value);
    const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(value);
    setPasswordValidation({
      minLength,
      hasUpper,
      hasLower,
      hasNumber,
      hasSpecial,
    });
  };

  const handleSendOTP = async (e) => {
    e.preventDefault();
    setError("");

    // --- CẢI TIẾN: Validate dữ liệu đầu vào ---
    if (verifyType === "email" && !/^\S+@\S+\.\S+$/.test(verifyKey)) {
      setError("Vui lòng nhập một địa chỉ email hợp lệ.");
      return;
    }
    if (verifyType === "phone" && !/^\d{10,11}$/.test(verifyKey)) {
      setError("Vui lòng nhập số điện thoại hợp lệ (10-11 chữ số).");
      return;
    }

    setLoading(true);
    try {
      await apiPost("/users/register", {
        verify_key: verifyKey.trim().toLowerCase(),
        verify_type: verifyType,
      });
      setStep(2);
    } catch (err) {
      // --- CẢI TIẾN: Hiển thị lỗi thân thiện hơn ---
      setError(
        err.response?.data?.error ||
          "Email hoặc SĐT này có thể đã được đăng ký."
      );
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOTP = async (e) => {
    e.preventDefault();
    setError("");

    // --- CẢI TIẾN: Validate OTP ---
    if (!/^\d{6}$/.test(otp)) {
      setError("Mã OTP phải là 6 chữ số.");
      return;
    }

    setLoading(true);
    try {
      await apiPost("/users/verify-otp", {
        verify_key: verifyKey.trim().toLowerCase(),
        otp,
      });
      setStep(3);
    } catch (err) {
      setError(
        err.response?.data?.error || "Mã OTP không hợp lệ hoặc đã hết hạn."
      );
    } finally {
      setLoading(false);
    }
  };

  const handleCompleteRegistration = async (e) => {
    e.preventDefault();
    setError("");

    // --- CẢI TIẾN: Validate các trường ---
    if (!names.trim()) {
      setError("Vui lòng nhập họ và tên của bạn.");
      return;
    }
    const isPasswordValid = Object.values(passwordValidation).every(Boolean);
    if (!isPasswordValid) {
      setError("Mật khẩu không đáp ứng các yêu cầu bảo mật.");
      return;
    }
    if (password !== confirmPassword) {
      setError(
        "Mật khẩu và xác nhận mật khẩu không khớp. Vui lòng kiểm tra lại."
      );
      return;
    }

    setLoading(true);
    try {
      await apiPost("/users/update-password-register", {
        verify_key: verifyKey.trim().toLowerCase(),
        user_password: password,
        user_names: names.trim(),
        user_type: userType,
        bio: bio.trim(),
      });
      // --- CẢI TIẾN: Chuyển hướng với thông báo ---
      navigate("/", {
        state: { message: "Đăng ký thành công!" },
      });
    } catch (err) {
      setError(
        err.response?.data?.error ||
          "Không thể hoàn tất đăng ký. Vui lòng thử lại."
      );
    } finally {
      setLoading(false);
    }
  };

  const renderStep = () => {
    switch (step) {
      case 1:
        return (
          <form onSubmit={handleSendOTP}>
            <p>Nhập email hoặc số điện thoại để bắt đầu.</p>
            <div className="verify-type-selector">
              <button
                type="button"
                className={verifyType === "email" ? "active" : ""}
                onClick={() => {
                  setVerifyType("email");
                  setVerifyKey("");
                }}
              >
                Email
              </button>
              <button
                type="button"
                className={verifyType === "phone" ? "active" : ""}
                onClick={() => {
                  setVerifyType("phone");
                  setVerifyKey("");
                }}
              >
                Số điện thoại
              </button>
            </div>
            <label htmlFor="verifyKey">
              {verifyType === "email" ? "Địa chỉ Email" : "Số điện thoại"}
            </label>
            <input
              type={verifyType === "email" ? "email" : "tel"}
              id="verifyKey"
              value={verifyKey}
              onChange={(e) => setVerifyKey(e.target.value)}
              required
              disabled={loading}
              placeholder={
                verifyType === "email" ? "your.email@example.com" : "0912345678"
              }
            />
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? "Đang gửi..." : "Tiếp tục"}
            </button>
          </form>
        );
      case 2:
        return (
          <form onSubmit={handleVerifyOTP}>
            <p>
              Chúng tôi đã gửi mã OTP đến <strong>{verifyKey}</strong>. Vui lòng
              nhập vào bên dưới.
            </p>
            <label htmlFor="otp">Mã xác thực</label>
            <input
              type="text"
              id="otp"
              value={otp}
              onChange={(e) => setOtp(e.target.value)}
              required
              disabled={loading}
              maxLength="6"
              placeholder="Mã gồm 6 chữ số"
            />
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? "Đang xác thực..." : "Xác thực"}
            </button>
            <button
              type="button"
              className="link-button"
              onClick={() => setStep(1)}
              disabled={loading}
            >
              Quay lại
            </button>
          </form>
        );
      case 3:
        return (
          <form onSubmit={handleCompleteRegistration}>
            <p>Hoàn tất hồ sơ của bạn để bắt đầu.</p>

            <label htmlFor="names">Họ và Tên</label>
            <input
              type="text"
              id="names"
              value={names}
              onChange={(e) => setNames(e.target.value)}
              required
              disabled={loading}
            />

            <label>Bạn là:</label>
            <div className="user-type-selection">
              <label>
                <input
                  type="radio"
                  name="userType"
                  value="buyer"
                  checked={userType === "buyer"}
                  onChange={() => setUserType("buyer")}
                />{" "}
                Người mua
              </label>
              <label>
                <input
                  type="radio"
                  name="userType"
                  value="seller"
                  checked={userType === "seller"}
                  onChange={() => setUserType("seller")}
                />{" "}
                Người bán
              </label>
            </div>

            <label htmlFor="bio">Tiểu sử ngắn</label>
            <textarea
              id="bio"
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Giới thiệu một chút về bạn..."
              disabled={loading}
              rows={3}
            />

            <label htmlFor="password">Mật khẩu</label>
            <div className="password-input-container">
              <input
                type={showPassword ? "text" : "password"}
                id="password"
                value={password}
                onChange={(e) => validatePassword(e.target.value)}
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
            {/* --- CẢI TIẾN: Hiển thị yêu cầu mật khẩu --- */}
            <div className="password-rules">
              <p className={passwordValidation.minLength ? "valid" : ""}>
                ✓ Ít nhất 8 ký tự
              </p>
              <p className={passwordValidation.hasUpper ? "valid" : ""}>
                ✓ Ít nhất 1 chữ hoa
              </p>
              <p className={passwordValidation.hasLower ? "valid" : ""}>
                ✓ Ít nhất 1 chữ thường
              </p>
              <p className={passwordValidation.hasNumber ? "valid" : ""}>
                ✓ Ít nhất 1 chữ số
              </p>
              <p className={passwordValidation.hasSpecial ? "valid" : ""}>
                ✓ Ít nhất 1 ký tự đặc biệt
              </p>
            </div>

            <label htmlFor="confirmPassword">Xác nhận Mật khẩu</label>
            <div className="password-input-container">
              <input
                type={showConfirmPassword ? "text" : "password"}
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                disabled={loading}
              />
              <span
                className="password-toggle-icon"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              >
                {showConfirmPassword ? (
                  <i className="fas fa-eye-slash"></i>
                ) : (
                  <i className="fas fa-eye"></i>
                )}
              </span>
            </div>
            {confirmPassword && password !== confirmPassword && (
              <p className="field-error-message">
                Mật khẩu xác nhận không khớp.
              </p>
            )}
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? "Đang tạo tài khoản..." : "Hoàn tất & Tạo tài khoản"}
            </button>
          </form>
        );
      default:
        return null;
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-form-container">
        <h1>Tham gia Taskly</h1>
        {error && <p className="error-message">{error}</p>}
        {renderStep()}
        <div className="form-footer">
          <span>Đã là thành viên? </span>
          <Link to="/login">Đăng nhập</Link>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
