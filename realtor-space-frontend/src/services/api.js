import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: 'https://real-estate-backend-840370620772.us-central1.run.app/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth endpoints
export const authApi = {
  login: (credentials) => api.post('/login', credentials),
  register: (userData) => api.post('/register', userData),
  refreshToken: () => api.post('/refresh'),
  logout: () => api.post('/logout'),
  requestPasswordReset: (email) => api.post('/auth/forgot-password', { email }),
  confirmPasswordReset: (token, password) => api.post('/auth/reset-password', { token, password }),
};

// Export individual auth functions for convenience
export const confirmPasswordReset = authApi.confirmPasswordReset;
export const requestPasswordReset = authApi.requestPasswordReset;

// Property endpoints  
export const propertyApi = {
  getAll: (params = {}) => api.get('/properties', { params }),
  getById: (id) => api.get(`/properties/${id}`),
  getMyProperties: (params = {}) => api.get('/my-properties', { params }),
  create: (propertyData) => api.post('/properties', propertyData),
  update: (id, propertyData) => api.put(`/properties/${id}`, propertyData),
  delete: (id) => api.delete(`/properties/${id}`),
  search: (searchParams) => api.get('/properties/search', { params: searchParams }),
  uploadImages: (id, formData) => api.post(`/properties/${id}/images`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
};

// User endpoints
export const userApi = {
  getProfile: () => api.get('/profile'),
  updateProfile: (userData) => api.put('/profile', userData),
  getUsers: () => api.get('/users'),
};

// Location endpoints
export const locationApi = {
  getCounties: () => api.get('/counties'),
  getSubCounties: (countyId) => api.get(`/counties/${countyId}/sub-counties`),
};

// For development, fallback to local JSON file
export const fallbackApi = {
  getProperties: async () => {
    try {
      const response = await fetch('/properties.json');
      const data = await response.json();
      return { data };
    } catch (error) {
      console.error('Error loading properties from fallback:', error);
      throw error;
    }
  }
};

export default api;
