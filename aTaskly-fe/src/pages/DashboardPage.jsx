import React, { useState, useEffect } from 'react';
import './DashboardPage.css';
import { apiGetAuth } from '../utils/api'; // Đảm bảo dùng apiGetAuth

const DashboardPage = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const fetchedOrders = await apiGetAuth('/user/orders'); // Đảm bảo dùng apiGetAuth
        setOrders(fetchedOrders);
      } catch (err) {
        const errorMessage = err.response?.data?.error || err.message;
        setError(errorMessage);
        console.error("Failed to fetch orders:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, []);

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="dashboard-page">
      <div className="dashboard-header">
        <h1>Dashboard</h1>
        <h2>Manage Orders</h2>
      </div>
      <div className="orders-table-container">
        <table className="orders-table">
          <thead>
            <tr>
              <th>Gig</th>
              <th>Buyer</th>
              <th>Price</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {orders.map(order => (
              <tr key={order._id}>
                <td className="gig-cell">
                  <img src={order.gig.images[0]} alt={order.gig.title} className="gig-thumbnail" />
                  <span>{order.gig.title}</span>
                </td>
                <td>{order.buyer.username}</td>
                <td>${order.price}</td>
                <td>
                  <span className={`status-badge status-${order.status.toLowerCase().replace(' ', '-')}`}>
                    {order.status}
                  </span>
                </td>
                <td>
                  <button className="action-button">Message</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DashboardPage; 