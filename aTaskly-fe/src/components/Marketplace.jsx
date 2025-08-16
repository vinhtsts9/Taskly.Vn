import React from 'react';
import { Link } from 'react-router-dom';
import './Marketplace.css';

const Marketplace = () => {
  const categories = [
    { name: 'Graphics & Design', slug: 'graphics-design' },
    { name: 'Digital Marketing', slug: 'digital-marketing' },
    { name: 'Writing & Translation', slug: 'writing-translation' },
    { name: 'Video & Animation', slug: 'video-animation' },
    { name: 'Music & Audio', slug: 'music-audio' },
    { name: 'Programming & Tech', slug: 'programming-tech' },
    { name: 'Business', slug: 'business' },
    { name: 'Lifestyle', slug: 'lifestyle' },
    { name: 'AI Services', slug: 'ai-services' },
  ];

  return (
    <section className="marketplace-section">
      <h2>Explore the marketplace</h2>
      <div className="category-grid">
        {categories.map((category) => (
          <Link to={`/category/${category.slug}`} className="category-item" key={category.name}>
            <div className="category-icon-placeholder"></div>
            <h3>{category.name}</h3>
          </Link>
        ))}
      </div>
    </section>
  );
};

export default Marketplace; 