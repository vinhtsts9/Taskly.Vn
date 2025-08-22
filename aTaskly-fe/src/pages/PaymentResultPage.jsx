import { useLocation } from "react-router-dom";
import { useEffect, useState } from "react";

const PaymentResultPage = () => {
  const { search } = useLocation();
  const params = new URLSearchParams(search);
  const success = params.get("success") === "true";
  const message = success ? "Thanh toán thành công" : "Thanh toán thất bại";
  const color = success ? "green" : "red";

  return (
    <div style={{ textAlign: "center", marginTop: 40 }}>
      <h2>Kết quả thanh toán</h2>
      <div
        style={{
          fontSize: 20,
          color: color,
        }}
      >
        {message}
      </div>
      {success === true && (
        <div style={{ marginTop: 20 }}>
          <a href="/" style={{ color: "#1890ff" }}>
            Quay về trang chủ
          </a>
        </div>
      )}
    </div>
  );
};

export default PaymentResultPage;
