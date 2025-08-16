import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Hero.css';

const Hero = () => {
  const navigate = useNavigate();
  const [keyword, setKeyword] = useState('');

  const suggestions = [
    'Thiết kế Logo',
    'Xây dựng Website',
    'Chạy quảng cáo',
    'Dịch thuật',
    'Kế toán',
    'Sửa lỗi phần mềm',
  ];

  const performSearch = (q) => {
    const query = (q ?? keyword).trim();
    if (query) {
      navigate(`/gigs?keyword=${encodeURIComponent(query)}`);
    } else {
      navigate('/gigs');
    }
  };

  const onKeyDown = (e) => {
    if (e.key === 'Enter') performSearch();
  };

  return (
    <section className="hero">
      <div className="hero__container">
        <div className="hero__eyebrow">Taskly.vn</div>
        <h1 className="hero__title">Kết nối công việc – hoàn thành trong ngày</h1>
        <p className="hero__subtitle">
          Thuê freelancer đáng tin cậy cho mọi đầu việc: thiết kế, lập trình, marketing, nội dung...
          Minh bạch – nhanh chóng – an tâm.
        </p>

        <div className="hero__search">
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onKeyDown={onKeyDown}
            placeholder="Tìm: 'thiết kế logo', 'sửa web WordPress'..."
          />
          <button onClick={() => performSearch()}>Tìm kiếm</button>
        </div>

        <div className="hero__chips">
          <span className="hero__chips-label">Gợi ý:</span>
          <ul>
            {suggestions.map((s) => (
              <li key={s} onClick={() => performSearch(s)}>{s}</li>
            ))}
          </ul>
        </div>

        <ul className="hero__features">
          <li>✓ Thanh toán an toàn</li>
          <li>✓ Chat và trao đổi tức thì</li>
          <li>✓ Bảo hành công việc</li>
        </ul>
      </div>
    </section>
  );
};

export default Hero;