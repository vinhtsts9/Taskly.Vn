import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './AuthForm.css';
import { apiPost } from '../utils/api'; // Sửa lại import

const RegisterPage = () => {
  const [step, setStep] = useState(1); // 1: Enter email/phone, 2: Verify OTP, 3: Set Password
  const [verifyKey, setVerifyKey] = useState('');
  const [verifyType, setVerifyType] = useState('email'); // 'email' or 'phone'
  const [otp, setOtp] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  // Thêm state cho các trường thông tin bổ sung
  const [names, setNames] = useState('');
  const [userType, setUserType] = useState('buyer'); // Mặc định là 'buyer'
  const [bio, setBio] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false); // Thêm trạng thái loading
  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      // Sửa lại payload để khớp với backend
      await apiPost('/users/register', { 
        verify_key: verifyKey, 
        verify_type: verifyType 
      });
      setStep(2);
    } catch (err) {
      // Hiển thị chỉ phần JSON data của response lỗi
      const rawResponse = err.response?.data ? JSON.stringify(err.response.data, null, 2) : err.message;
      setError(rawResponse);
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOTP = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await apiPost('/users/verify-otp', { verify_key: verifyKey, 
        otp: otp });
      setStep(3);
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message;
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleCompleteRegistration = async (e) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setError('Mật khẩu không khớp.');
      return;
    }
    setError('');
    setLoading(true);
    try {
      // Gửi đầy đủ thông tin để hoàn tất đăng ký
      await apiPost('/users/complete-registration', { 
        verify_key: verifyKey,
        user_password: password,
        user_names: names,
        user_type: userType,
        bio: bio
      });
      alert('Đăng ký thành công! Bạn sẽ được chuyển đến trang đăng nhập.');
      navigate('/login');
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message;
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const renderStep = () => {
    switch (step) {
      case 1:
        return (
          <form onSubmit={handleRegister}>
            <p>Enter your email or phone number to start.</p>
            <div className="verify-type-selector">
              <button type="button" className={verifyType === 'email' ? 'active' : ''} onClick={() => setVerifyType('email')}>Email</button>
              <button type="button" className={verifyType === 'phone' ? 'active' : ''} onClick={() => setVerifyType('phone')}>Phone</button>
            </div>
            <label htmlFor="verifyKey">{verifyType === 'email' ? 'Email' : 'Phone Number'}</label>
            <input
              type={verifyType === 'email' ? 'email' : 'tel'}
              id="verifyKey"
              value={verifyKey}
              onChange={(e) => setVerifyKey(e.target.value)}
              required
              disabled={loading}
            />
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? 'Sending...' : 'Continue'}
            </button>
          </form>
        );
      case 2:
        return (
          <form onSubmit={handleVerifyOTP}>
            <p>We've sent an OTP to {verifyKey}. Please enter it below.</p>
            <label htmlFor="otp">Verification Code</label>
            <input
              type="text"
              id="otp"
              value={otp}
              onChange={(e) => setOtp(e.target.value)}
              required
              disabled={loading}
            />
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? 'Verifying...' : 'Verify'}
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

            <label>Loại tài khoản</label>
            <div className="verify-type-selector">
              <button type="button" className={userType === 'buyer' ? 'active' : ''} onClick={() => setUserType('buyer')}>Tôi là người mua</button>
              <button type="button" className={userType === 'seller' ? 'active' : ''} onClick={() => setUserType('seller')}>Tôi là người bán</button>
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
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={loading}
            />
            <label htmlFor="confirmPassword">Xác nhận Mật khẩu</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              disabled={loading}
            />
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? 'Đang tạo tài khoản...' : 'Hoàn tất & Tạo tài khoản'}
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
        {error && <pre className="error-message">{error}</pre>}
        {renderStep()}
        <div className="form-footer">
          <span>Already a member? </span>
          <Link to="/login">Sign In</Link>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage; 