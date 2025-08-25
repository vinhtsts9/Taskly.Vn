const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

// Biến cờ để ngăn chặn nhiều yêu cầu làm mới token cùng lúc
let isRefreshing = false;
// Hàng đợi cho các yêu cầu bị lỗi 401 trong khi đang làm mới token
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

const apiRequest = async (endpoint, method, data = null, options = {}) => {
  const fullUrl = `${API_BASE_URL}${endpoint}`;
  const isForm = typeof FormData !== "undefined" && data instanceof FormData;
  const { needsAuth = false, idempotencyKey = null } = options;
  const headers = {};
  const fetchOptions = { method, headers };

  console.log("env", import.meta.env);
  if (idempotencyKey) {
    headers["Idempotency-Key"] = idempotencyKey;
  }

  if (data) {
    if (isForm) {
      fetchOptions.body = data; // Sửa ở đây
    } else {
      headers["Content-Type"] = "application/json";
      fetchOptions.body = JSON.stringify(data); // Và ở đây
    }
  }

  // Chỉ gửi credentials và xử lý refresh token cho các yêu cầu cần xác thực
  if (needsAuth) {
    fetchOptions.credentials = "include";
  }

  try {
    let response = await fetch(fullUrl, fetchOptions);

    // Logic làm mới token CHỈ áp dụng cho các yêu cầu xác thực
    if (needsAuth && response.status === 401) {
      if (!isRefreshing) {
        isRefreshing = true;
        try {
          const refreshResponse = await fetch(
            `${API_BASE_URL}/users/refresh-token`,
            {
              method: "POST",
              credentials: "include",
            }
          );

          if (!refreshResponse.ok) {
            const err = new Error("Phiên đăng nhập đã hết hạn.");
            processQueue(err, null);
            window.dispatchEvent(new Event("auth-failure"));
            throw err;
          }

          processQueue(null, null);
          // Thử lại yêu cầu gốc với cùng options (đã có credentials: 'include')
          response = await fetch(fullUrl, fetchOptions);
        } catch (error) {
          processQueue(error, null);
          throw error;
        } finally {
          isRefreshing = false;
        }
      } else {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => fetch(fullUrl, fetchOptions))
          .then((res) => {
            if (!res.ok)
              throw new Error(
                `Yêu cầu thất bại sau khi làm mới token với mã trạng thái ${res.status}`
              );
            const contentType = res.headers.get("content-type");
            if (contentType && contentType.includes("application/json")) {
              return res.json();
            }
            return { success: true };
          });
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const error = new Error(
        errorData.error ||
          errorData.message ||
          `Yêu cầu thất bại với mã trạng thái ${response.status}`
      );
      error.response = { data: errorData, status: response.status };
      throw error;
    }

    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      return response.json();
    }

    return { success: true };
  } catch (error) {
    console.error("Lỗi API:", error);
    if (needsAuth && (!error.response || error.response.status === 401)) {
      window.dispatchEvent(new Event("auth-failure"));
    }
    throw error;
  }
};

// Các hàm public - không gửi credentials
export const apiGet = (endpoint) =>
  apiRequest(endpoint, "GET", null, { needsAuth: false });
export const apiPost = (endpoint, data) =>
  apiRequest(endpoint, "POST", data, { needsAuth: false });

// Các hàm cần xác thực - sẽ gửi credentials và có logic refresh token
export const apiGetAuth = (endpoint, options) =>
  apiRequest(endpoint, "GET", null, { ...options, needsAuth: true });
export const apiPostAuth = (endpoint, data, options) =>
  apiRequest(endpoint, "POST", data, { ...options, needsAuth: true });
export const apiPutAuth = (endpoint, data, options) =>
  apiRequest(endpoint, "PUT", data, { ...options, needsAuth: true });
export const apiDeleteAuth = (endpoint, options) =>
  apiRequest(endpoint, "DELETE", null, { ...options, needsAuth: true });
// Một số API DELETE của backend nhận body JSON
export const apiDeleteAuthWithBody = (endpoint, data, options) =>
  apiRequest(endpoint, "DELETE", data, { ...options, needsAuth: true });
// Upload file (multipart/form-data)
export const apiUploadAuth = (endpoint, formData, options) =>
  apiRequest(endpoint, "POST", formData, { ...options, needsAuth: true });
