// Helpers to evaluate user roles without tightly coupling to backend shape

export const isUserAdmin = (user) => {
  if (!user) return false;

  // Common flags
  if (user.isAdmin === true) return true;

  // user_type as array of strings
  if (Array.isArray(user.user_type)) {
    const hasAdmin = user.user_type.some((t) => typeof t === 'string' && t.toLowerCase() === 'admin');
    if (hasAdmin) return true;
  }

  return false;
};

