import React from 'react';
import { Link } from 'react-router-dom';
import './GigCard.css';

const GigCard = ({ gig }) => {
  return (
    <Link to={`/gig/${gig.id}`} className="gig-card-link">
      <div className="gig-card">
        <img 
          src={gig.image_url && gig.image_url.length > 0 ? gig.image_url[0] : 'https://via.placeholder.com/300'}
          alt={gig.title}
          className="gig-image"
        />
        <div className="gig-seller-info">
          <div className="seller-avatar-placeholder"></div>
          {/* <span>{gig.seller.name}</span> */}
        </div>
        <p className="gig-title">{gig.title}</p>
        <div className="gig-rating">
          <span>⭐ {gig.rating}</span>
          <span>({gig.reviews})</span>
        </div>
        <div className="gig-price">
          <span>STARTING AT</span>
          <strong>${gig.basic_price}</strong>
        </div>
      </div>
    </Link>
  );
};

export default GigCard; 