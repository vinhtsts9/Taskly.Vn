import React from 'react';
import './ServiceCard.css';

const ServiceCard = ({ title, category }) => {
  return (
    <div className="service-card">
      <div className="card-image-placeholder"></div>
      <div className="card-content">
        <h3>{category}</h3>
        <h2>{title}</h2>
      </div>
    </div>
  );
};

export default ServiceCard; 