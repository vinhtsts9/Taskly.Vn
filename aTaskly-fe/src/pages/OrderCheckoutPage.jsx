import React, { useState, useContext, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import { apiPostAuth } from "../utils/api";
import { AuthContext } from "../context/AuthContext";
import "./OrderCheckoutPage.css";

const OrderCheckoutPage = () => {
  const { state } = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated } = useContext(AuthContext);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [isValid, setIsValid] = useState(false);

  // Chuyển hướng nếu không có dữ liệu đơn hàng hoặc người dùng chưa đăng nhập
  useEffect(() => {
    // Thực hiện kiểm tra và điều hướng trong useEffect
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }
    if (!state || !state.order || !state.gig) {
      navigate("/"); // Quay về trang chủ nếu thiếu dữ liệu
      return;
    }
    // Nếu tất cả đều hợp lệ, cho phép component render
    setIsValid(true);
  }, [isAuthenticated, state, navigate]);

  // Chỉ render nội dung khi đã xác thực và có đủ dữ liệu
  if (!isValid) {
    return null; // Hoặc một component loading
  }

  const { order, gig } = state;

  const formatPrice = (price) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  const handlePayment = async () => {
    setLoading(true);
    setError("");

    try {
      const idempotencyKey = uuidv4();
      const payload = {
        order_id: order.id,
        payment_method: "vnpay",
      };

      // Gọi API backend để lấy URL thanh toán của VNPay
      const response = await apiPostAuth("/payments/create-intent", payload, {
        idempotencyKey,
      });
      console.log("response payment", response);
      if (response && response.payment_url) {
        // Chuyển hướng người dùng đến cổng thanh toán của VNPay
        window.location.href = response.payment_url.payment_url;
      } else {
        setError("Không nhận được URL thanh toán. Vui lòng thử lại.");
      }
    } catch (err) {
      const errorMessage =
        err.response?.data?.error ||
        err.message ||
        "Đã xảy ra lỗi khi khởi tạo thanh toán.";
      console.log("error payment", errorMessage);
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="order-checkout-page">
      <h1>Thanh toán đơn hàng</h1>
      <div className="order-summary-card">
        <h2>Tóm tắt đơn hàng</h2>
        <div className="order-detail-item">
          <span>Dịch vụ:</span>
          <span>{gig.title}</span>
        </div>
        <div className="order-detail-item">
          <span>Gói:</span>
          <span className="package-tier">{order.package_tier}</span>
        </div>
        <hr />
        <div className="order-detail-item total-price">
          <span>Tổng cộng:</span>
          <span>{formatPrice(order.total_price)}</span>
        </div>
      </div>

      {error && <p className="error-message">{error}</p>}

      <div className="payment-actions">
        <button
          className="vnpay-button"
          onClick={handlePayment}
          disabled={loading}
        >
          {loading ? "Đang xử lý..." : "Thanh toán qua VNPay"}
        </button>
        <p className="payment-note">
          Bạn sẽ được chuyển đến cổng thanh toán an toàn của VNPay.
        </p>
      </div>
    </div>
  );
};

export default OrderCheckoutPage;
