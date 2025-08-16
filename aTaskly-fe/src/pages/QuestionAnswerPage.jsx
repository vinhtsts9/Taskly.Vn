import React, { useState, useEffect, useContext } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "./QuestionAnswerPage.css";
import { AuthContext } from "../context/AuthContext";
import { apiPostAuth } from "../utils/api";

const QuestionAnswerPage = () => {
  const { state } = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated, currentUser } = useContext(AuthContext);

  const [gig, setGig] = useState(null);
  const [selectedPackage, setSelectedPackage] = useState(null);
  const [answers, setAnswers] = useState({});

  // State cho modal và trạng thái gọi API
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    if (!state || !state.gig || !state.selectedPackage) {
      navigate("/404");
      return;
    }
    console.log(" state:", state);

    const { gig: gigData, selectedPackage: pkg } = state;
    setGig(gigData);
    setSelectedPackage(pkg);

    if (gigData.Question && gigData.Question.length > 0) {
      const initialAnswers = gigData.Question.reduce((acc, question) => {
        acc[question.id] = ""; // Sửa: Dùng question.ID (viết hoa) để khớp với API
        return acc;
      }, {});
      setAnswers(initialAnswers);
    }
  }, [state, navigate]);

  const formatPrice = (price) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  // Bước 1: Kiểm tra dữ liệu và mở hộp thoại xác nhận
  const handleOpenConfirmationModal = () => {
    setError(""); // Xóa lỗi cũ
    setSuccessMessage("");

    if (!isAuthenticated) {
      navigate("/login");
      return;
    }

    if (!gig || !selectedPackage) {
      setError("Thông tin gig hoặc gói dịch vụ không đầy đủ.");
      return;
    }

    // Kiểm tra các câu hỏi bắt buộc
    if (gig.Question && gig.Question.length > 0) {
      const missingRequiredAnswers = gig.Question.some(
        (q) => q.required && !answers[q.id]?.trim()
      );
      if (missingRequiredAnswers) {
        setError("Vui lòng trả lời tất cả các câu hỏi bắt buộc.");
        console.error("Missing required answers:", answers);
        return;
      }
    }

    // Nếu mọi thứ hợp lệ, mở hộp thoại
    setIsModalOpen(true);
  };

  // Bước 2: Tạo đơn hàng khi người dùng bấm "Xác nhận" trong hộp thoại
  const handleCreateOrder = async () => {
    setLoading(true);
    setError("");
    setSuccessMessage("");

    try {
      const answersPayload = Object.keys(answers).map((questionId) => ({
        question_id: questionId,
        answer: answers[questionId], // giữ nguyên, không trim
      }));

      const deliveryDate = new Date();
      deliveryDate.setDate(
        deliveryDate.getDate() + selectedPackage.delivery_days
      );

      const orderPayload = {
        // Payload an toàn: Backend sẽ tự lấy buyer_id và tính total_price
        gig_id: gig.id,
        seller_id: gig.user_id,
        package_tier: selectedPackage.tier,
        delivery_date: deliveryDate.toISOString(),
        answers: answersPayload,
      };
      console.log("Order Payload:", orderPayload);
      const newOrder = await apiPostAuth("/orders/create", orderPayload);

      // Thành công: hiển thị thông báo trong hộp thoại
      setSuccessMessage("Đơn hàng đã được tạo thành công!");

      // Thêm lại logic: Tự động chuyển hướng đến trang thanh toán sau 2 giây
      setTimeout(() => {
        setIsModalOpen(false);
        // API trả về { message: "...", order: {...} }
        navigate("/order-checkout", { state: { order: newOrder.order } });
      }, 2000);
    } catch (err) {
      const errorMessage =
        err.response?.data?.error ||
        err.message ||
        "Đã xảy ra lỗi khi tạo đơn hàng.";
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  if (!gig || !selectedPackage) {
    return <div className="not-found">Không tìm thấy thông tin cần thiết.</div>;
  }

  return (
    <div className="question-answer-page">
      <h1>Trả lời câu hỏi cho Gig: {gig.title}</h1>
      <div className="selected-package-summary">
        <h2>Gói đã chọn: {selectedPackage.tier}</h2>
        <p>Giá: {formatPrice(selectedPackage.price)}</p>
        <p>Giao hàng trong: {selectedPackage.delivery_days} ngày</p>
      </div>

      {error && !isModalOpen && <p className="error-message">{error}</p>}

      {gig.Question && gig.Question.length > 0 ? (
        <div className="questions-section" key={gig.Question[0].id}>
          <h2>Các câu hỏi của người bán</h2>
          {gig.Question.map((question) => (
            <div key={question.id} className="question-item">
              <label className="question-text">
                {question.question}
                {question.required && (
                  <span className="required-badge"> (Bắt buộc)</span>
                )}
              </label>
              <textarea
                className="question-answer-input"
                placeholder="Nhập câu trả lời của bạn..."
                value={answers[question.id] || ""}
                onChange={(e) =>
                  setAnswers({ ...answers, [question.id]: e.target.value })
                }
                rows="4"
              ></textarea>
            </div>
          ))}
        </div>
      ) : (
        <div className="no-questions">
          <p>
            Không có câu hỏi nào từ người bán cho gig này. Bạn có thể tiếp tục.
          </p>
        </div>
      )}

      <button
        className="continue-to-payment-button"
        onClick={handleOpenConfirmationModal}
      >
        Tiếp tục
      </button>

      {/* Hộp thoại xác nhận */}
      {isModalOpen && (
        <div className="modal-overlay">
          <div className="modal-content">
            {!successMessage ? (
              <>
                <h2>Xác nhận đơn hàng</h2>
                <p>Bạn sắp tạo một đơn hàng mới với tổng giá trị là:</p>
                <p className="modal-price">
                  {formatPrice(selectedPackage.price)}
                </p>

                {error && <p className="error-message modal-error">{error}</p>}

                <div className="modal-actions">
                  <button
                    className="modal-button cancel"
                    onClick={() => setIsModalOpen(false)}
                    disabled={loading}
                  >
                    Hủy
                  </button>
                  <button
                    className="modal-button confirm"
                    onClick={handleCreateOrder}
                    disabled={loading}
                  >
                    {loading ? "Đang xử lý..." : "Xác nhận & Tạo đơn"}
                  </button>
                </div>
              </>
            ) : (
              <div className="success-message-container">
                <svg
                  className="success-icon"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 52 52"
                >
                  <circle
                    className="success-icon-circle"
                    cx="26"
                    cy="26"
                    r="25"
                    fill="none"
                  />
                  <path
                    className="success-icon-checkmark"
                    fill="none"
                    d="M14.1 27.2l7.1 7.2 16.7-16.8"
                  />
                </svg>
                <h2>Thành công!</h2>
                <p>{successMessage}</p>
                <p>Đang chuyển hướng đến trang thanh toán...</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default QuestionAnswerPage;
