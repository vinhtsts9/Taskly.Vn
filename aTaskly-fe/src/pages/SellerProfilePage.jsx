import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import GigCard from '../components/GigCard';
import './SellerProfilePage.css';
import { apiGet } from '../utils/api';

const SellerProfilePage = () => {
  const { userId } = useParams();
  const [seller, setSeller] = useState(null);
  const [gigs, setGigs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchSellerData = async () => {
      setLoading(true);
      try {
        const [sellerData, gigsData] = await Promise.all([
          apiGet(`/users/${userId}`),
          apiGet(`/users/${userId}/gigs`)
        ]);
        setSeller(sellerData);
        setGigs(gigsData);
      } catch (err) {
        const errorMessage = err.response?.data?.error || err.message;
        setError(errorMessage);
        console.error("Failed to fetch seller data:", err);
      } finally {
        setLoading(false);
      }
    };

    if (userId) {
      fetchSellerData();
    }
  }, [userId]);

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Error: {error}</div>;
  if (!seller) return <div className="not-found">Seller not found.</div>;

  return (
    <div className="seller-profile-page">
      <div className="profile-sidebar">
        <div className="profile-card">
          <img src={seller.avatar} alt={seller.username} className="profile-avatar" />
          <h2 className="profile-name">{seller.username}</h2>
          <p className="profile-level">{seller.level || 'New Seller'}</p>
          <hr />
          <ul className="profile-stats">
            <li><span>From</span><strong>{seller.country}</strong></li>
            <li><span>Member since</span><strong>{new Date(seller.createdAt).toLocaleDateString()}</strong></li>
            <li><span>Avg. response time</span><strong>{seller.avgResponseTime || 'N/A'}</strong></li>
          </ul>
        </div>
        <div className="profile-description-card">
          <h3>Description</h3>
          <p>{seller.description}</p>
          <h3>Skills</h3>
          <div className="skills-container">
            {seller.skills.map(skill => <span key={skill} className="skill-tag">{skill}</span>)}
          </div>
        </div>
      </div>
      <div className="gigs-main-content">
        <h2>{seller.username}'s Gigs</h2>
        <div className="gigs-grid">
          {gigs.map(gig => (
            <GigCard key={gig.id} gig={gig} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default SellerProfilePage; 