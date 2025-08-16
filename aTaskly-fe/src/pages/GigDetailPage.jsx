import React, { useState, useEffect, useContext } from "react";
import { useParams, useNavigate } from "react-router-dom";
import "./GigDetailPage.css";
import { apiGet, apiPostAuth, apiGetAuth } from "../utils/api";
import { AuthContext } from "../context/AuthContext";

const GigDetailPage = () => {
  const { gigId } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated, currentUser } = useContext(AuthContext);
  const [gig, setGig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedPackage, setSelectedPackage] = useState(null);
  const [contactLoading, setContactLoading] = useState(false);

  const handleContinue = () => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }

    if (selectedPackage) {
      navigate("/gig-questions", {
        state: {
          gig: gig,
          selectedPackage: selectedPackage,
        },
      });
    }
  };

  useEffect(() => {
    const fetchGigDetail = async () => {
      setLoading(true);
      try {
        const data = await apiGet(`/gigs/${gigId}`);
        setGig(data);
        // Tự động chọn package đầu tiên nếu có
        if (data.GigPackage && data.GigPackage.length > 0) {
          setSelectedPackage(data.GigPackage[0]);
        }
      } catch (err) {
        const errorMessage = err.response?.data?.error || err.message;
        setError(errorMessage);
        console.error("Failed to fetch gig details:", err);
      } finally {
        setLoading(false);
      }
    };

    if (gigId) {
      fetchGigDetail();
    }
  }, [gigId]);

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Error: {error}</div>;
  if (!gig) return <div className="not-found">Gig not found.</div>;

  const formatPrice = (price) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  const getTierDisplayName = (tier) => {
    const tierNames = {
      basic: "basic",
      standard: "standard",
      premium: "premium",
    };
    return tierNames[tier] || tier;
  };

  const handleContactSeller = async () => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }

    if (!currentUser || !gig) {
      return;
    }

    setContactLoading(true);
    try {
      // 1) Kiểm tra xem đã có phòng giữa 2 người chưa
      const payload = {
        user2_id: gig.user_id,
      };
      const roomExists = await apiPostAuth("/chat/room-exists", payload);

      // Lưu seller info để trang chat sử dụng
      const sellerInfo = {
        id: gig.user_id,
        name: gig.user_name,
        avatar: gig.user_profile_pic || "",
      };
      localStorage.setItem("temp_seller_info", JSON.stringify(sellerInfo));

      if (
        roomExists &&
        roomExists.id &&
        roomExists.id !== "00000000-0000-0000-0000-000000000000"
      ) {
        // 2) Nếu đã có phòng -> điều hướng trực tiếp vào phòng thật (WS + history)
        navigate(`/messages?room=${roomExists.id}`);
      } else {
        // 3) Nếu chưa có -> vào phòng tạm như luồng hiện tại
        navigate("/messages?temp_room=true");
      }
    } catch (err) {
      // Nếu lỗi, fallback vào phòng tạm để vẫn cho phép người dùng nhắn
      navigate("/messages?temp_room=true");
    } finally {
      setContactLoading(false);
    }
  };

  return (
    <div className="gig-detail-page">
      <div className="main-content">
        <div className="breadcrumb">
          {gig.category_name ? gig.category_name : "CATEGORY"}
        </div>
        <h1 className="gig-title">{gig.title}</h1>

        <div className="seller-info">
          <img
            src={gig.user_profile_pic || "/default-avatar.png"}
            alt={gig.user_name}
            className="seller-avatar"
          />
          <span className="seller-name">{gig.user_name}</span>
        </div>

        <div className="gig-gallery">
          {gig.image_url && gig.image_url.length > 0 && (
            <div className="gallery-container">
              <img
                src={gig.image_url[0]}
                alt="Gig main image"
                className="main-image"
              />
              {gig.image_url.length > 1 && (
                <div className="thumbnail-images">
                  {gig.image_url.slice(1).map((url, index) => (
                    <img
                      key={index}
                      src={url}
                      alt={`Gig image ${index + 2}`}
                      className="thumbnail"
                    />
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        <div className="gig-description">
          <h2>Về dịch vụ này</h2>
          <div dangerouslySetInnerHTML={{ __html: gig.description }} />
        </div>

        {/* Packages Section */}
        {gig.GigPackage && gig.GigPackage.length > 0 && (
          <div className="gig-packages">
            <h2>Gói dịch vụ</h2>
            <div className="packages-grid">
              {gig.GigPackage.map((pkg, index) => (
                <div
                  key={index}
                  className={`package-card ${
                    selectedPackage === pkg ? "selected" : ""
                  }`}
                  onClick={() => setSelectedPackage(pkg)}
                >
                  <div className="package-header">
                    <h3 className="package-tier">
                      {getTierDisplayName(pkg.tier)}
                    </h3>
                    <div className="package-price">
                      {formatPrice(pkg.price)}
                    </div>
                  </div>
                  <div className="package-delivery">
                    <span>Giao hàng trong {pkg.delivery_days} ngày</span>
                  </div>
                  {pkg.options && (
                    <div className="package-options">
                      <div className="option-item">
                        <span className="option-label">Số lần sửa:</span>
                        <span className="option-value">
                          {pkg.options.revisions || 0}
                        </span>
                      </div>
                      <div className="option-item">
                        <span className="option-label">Số file bàn giao:</span>
                        <span className="option-value">
                          {pkg.options.files || 1}
                        </span>
                      </div>
                    </div>
                  )}
                  <button className="select-package-btn">Chọn gói này</button>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Requirements Section */}
        {gig.Question && gig.Question.length > 0 && (
          <div className="gig-requirements">
            <h2>Câu hỏi cho người mua</h2>
            <div className="requirements-list">
              {gig.Question.map((question, index) => (
                <div key={index} className="requirement-item">
                  <div className="question-header">
                    <span className="question-text">{question.question}</span>
                    {question.required && (
                      <span className="required-badge">Bắt buộc</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="about-seller">
          <h2>Về người bán</h2>
          <p>
            Thông tin chi tiết về {gig.user_name} sẽ có sẵn trong tương lai.
          </p>
        </div>

        <div className="gig-reviews">
          <h2>Đánh giá</h2>
          <p>Đánh giá sẽ được hiển thị ở đây trong tương lai.</p>
        </div>
      </div>

      <div className="sidebar">
        <div className="pricing-box">
          <div className="pricing-header">
            <h3>{gig.title}</h3>
            {selectedPackage ? (
              <span className="price">
                {formatPrice(selectedPackage.price)}
              </span>
            ) : (
              <span className="price">Chọn gói</span>
            )}
          </div>

          {selectedPackage && (
            <>
              <p className="delivery-info">
                Giao hàng trong {selectedPackage.delivery_days} ngày
              </p>
              {selectedPackage.options && (
                <div className="selected-package-options">
                  <div className="option-summary">
                    <span>
                      Sửa {selectedPackage.options.revisions || 0} lần
                    </span>
                    <span>
                      Bàn giao {selectedPackage.options.files || 1} file
                    </span>
                  </div>
                </div>
              )}
            </>
          )}

          <button
            className="continue-button"
            onClick={handleContinue}
            disabled={!selectedPackage}
          >
            {selectedPackage
              ? `Tiếp tục (${formatPrice(selectedPackage.price)})`
              : "Chọn gói để tiếp tục"}
          </button>

          <button
            className="contact-seller"
            onClick={handleContactSeller}
            disabled={contactLoading}
          >
            {contactLoading ? "Đang xử lý..." : "Liên hệ người bán"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default GigDetailPage;
