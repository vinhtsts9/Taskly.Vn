import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import RequireAdmin from './components/RequireAdmin';
import HomePage from './pages/HomePage';
import CategoryPage from './pages/CategoryPage';
import GigDetailPage from './pages/GigDetailPage';
import SellerProfilePage from './pages/SellerProfilePage';
import DashboardPage from './pages/DashboardPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import GigsPage from './pages/GigsPage'; // Import GigsPage
import ChatPage from './pages/ChatPage'; // Import ChatPage
import AdminDashboardPage from './pages/AdminDashboardPage';
import GigCreatePage from './pages/GigCreatePage';
import MyGigsPage from './pages/MyGigsPage'; // Import MyGigsPage
import QuestionAnswerPage from './pages/QuestionAnswerPage'; // Import QuestionAnswerPage

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Routes with main layout */}
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="gigs" element={<GigsPage />} /> {/* Add GigsPage route */}
          <Route path="category/:categoryName" element={<CategoryPage />} />
          <Route path="gig/:gigId" element={<GigDetailPage />} />
          <Route path="seller/:sellerId" element={<SellerProfilePage />} />
          <Route path="dashboard/orders" element={<DashboardPage />} />
          <Route path="messages" element={<ChatPage />} /> {/* Add ChatPage route */}
          <Route
            path="admin"
            element={
              <RequireAdmin>
                <AdminDashboardPage />
              </RequireAdmin>
            }
          />
          <Route path="gigs/new" element={<GigCreatePage />} />
          <Route path="my-gigs" element={<MyGigsPage />} /> {/* Add MyGigsPage route */}
          <Route path="gig-questions" element={<QuestionAnswerPage />} /> {/* Add QuestionAnswerPage route */}
        </Route>

        {/* Auth routes without main layout */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App; 