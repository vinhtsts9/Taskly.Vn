import React from 'react';
import './Footer.css';

const Footer = () => {
  const footerColumns = [
    {
      title: 'Danh mục',
      links: ['Thiết kế & Đồ họa', 'Marketing Online', 'Viết & Dịch thuật', 'Video & Hoạt họa'],
    },
    {
      title: 'Về Taskly',
      links: ['Nghề nghiệp', 'Báo chí', 'Đối tác', 'Chính sách riêng tư'],
    },
    {
      title: 'Hỗ trợ',
      links: ['Trung tâm hỗ trợ', 'Tin cậy & An toàn', 'Bán hàng trên Taskly', 'Mua hàng trên Taskly'],
    },
    {
      title: 'Cộng đồng',
      links: ['Sự kiện', 'Blog', 'Diễn đàn', 'Podcast'],
    },
  ];

  return (
    <footer className="site-footer">
      <div className="footer-columns">
        {footerColumns.map((column) => (
          <div className="footer-column" key={column.title}>
            <h3>{column.title}</h3>
            <ul>
              {column.links.map((link) => (
                <li key={link}>
                  <a href="#">{link}</a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
      <div className="footer-bottom">
        <div className="footer-logo">
          <span>Taskly</span>
          <p>© Taskly Việt Nam 2024</p>
        </div>
        <div className="footer-social">
          {/* Placeholder for social icons */}
          <span>FB TW IN LI</span>
        </div>
      </div>
    </footer>
  );
};

export default Footer; 