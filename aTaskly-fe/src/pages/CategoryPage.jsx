import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import GigCard from '../components/GigCard';
import { apiGet } from '../utils/api'; // Import apiGet
import './CategoryPage.css';

const CategoryPage = () => {
  const { categoryName } = useParams();
  const [gigs, setGigs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchGigs = async () => {
      setLoading(true);
      setError('');
      try {
        const data = await apiGet(`/user/gigs?category=${categoryName}`);
        setGigs(data.gigs || []);
      } catch (err) {
        const errorMessage = err.response?.data?.error || 'Failed to fetch gigs for this category.';
        setError(errorMessage);
      } finally {
        setLoading(false);
      }
    };

    fetchGigs();
  }, [categoryName]); // Chạy lại mỗi khi categoryName trên URL thay đổi

  const formattedCategoryName = categoryName
    .replace(/-/g, ' ')
    .replace(/\b\w/g, (l) => l.toUpperCase());

  if (loading) {
    return <div className="loading-state">Loading services...</div>;
  }

  if (error) {
    return <div className="error-message">{error}</div>;
  }

  return (
    <div className="category-page">
      <div className="category-header">
        <h1>{formattedCategoryName}</h1>
        <p>Explore the boundaries of art and technology with Fiverr's talented designers.</p>
      </div>
      <div className="gig-list">
        {gigs.length > 0 ? (
          gigs.map(gig => <GigCard key={gig._id} gig={gig} />)
        ) : (
          <p>No services found in this category.</p>
        )}
      </div>
    </div>
  );
};

export default CategoryPage; 