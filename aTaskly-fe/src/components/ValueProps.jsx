import React from 'react';
import './ValueProps.css';

const ValueProps = () => {
  const points = [
    {
      title: 'Stick to your budget',
      description: 'Find the right service for every price point. No hourly rates, just project-based pricing.',
    },
    {
      title: 'Get quality work done quickly',
      description: 'Hand your project over to a talented freelancer in minutes, get long-lasting results.',
    },
    {
      title: 'Pay when you\'re happy',
      description: 'Upfront quotes mean no surprises. Payments are only released when you approve the work.',
    },
    {
      title: 'Count on 24/7 support',
      description: 'Our round-the-clock support team is available to help anytime, anywhere.',
    },
  ];

  return (
    <section className="value-props-section">
      <div className="value-props-content">
        <h2>The best part? Everything.</h2>
        <ul>
          {points.map((point) => (
            <li key={point.title}>
              <h3>{point.title}</h3>
              <p>{point.description}</p>
            </li>
          ))}
        </ul>
      </div>
      <div className="value-props-media">
        {/* Placeholder for an image or video */}
        <div className="media-placeholder"></div>
      </div>
    </section>
  );
};

export default ValueProps; 