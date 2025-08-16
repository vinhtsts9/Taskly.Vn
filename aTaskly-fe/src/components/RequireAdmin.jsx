import React, { useContext } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { isUserAdmin } from '../utils/auth';

const RequireAdmin = ({ children }) => {
  const { currentUser, isAuthenticated } = useContext(AuthContext);
  const location = useLocation();

  // Debug logs to verify auth state and role detection
  try {
    // eslint-disable-next-line no-console
    console.debug('[RequireAdmin] isAuthenticated:', isAuthenticated, 'isAdmin:', isUserAdmin(currentUser), 'currentUser:', currentUser, 'from:', location.pathname);
  } catch (_) {}

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  if (!isUserAdmin(currentUser)) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default RequireAdmin;



