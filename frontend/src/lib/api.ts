import axios from 'axios';

// Create an Axios instance
const api = axios.create({
  baseURL: '/api', // Nginx will proxy this to the backend
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Important for cookies (if using HttpOnly cookies for refresh tokens)
});

// Request Interceptor (Attach Access Token)
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken'); // Or from memory/context
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response Interceptor (Handle Token Refresh)
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If 401 and not already retrying
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Call refresh endpoint (assumes HttpOnly cookie for refresh token)
        const { data } = await axios.post('/api/auth/refresh', {}, { withCredentials: true });
        
        // Update local storage or state with new access token
        localStorage.setItem('accessToken', data.accessToken);
        
        // Update header for the original request
        originalRequest.headers.Authorization = `Bearer ${data.accessToken}`;
        
        // Retry original request
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed - redirect to login
        localStorage.removeItem('accessToken');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
