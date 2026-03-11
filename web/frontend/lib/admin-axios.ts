import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const ADMIN_TOKEN_STORAGE_KEY = "admin_token";

const adminAxiosClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

adminAxiosClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (typeof window !== "undefined") {
      const token = localStorage.getItem(ADMIN_TOKEN_STORAGE_KEY);
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

adminAxiosClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401 && typeof window !== "undefined") {
      localStorage.removeItem(ADMIN_TOKEN_STORAGE_KEY);
      window.dispatchEvent(new Event("admin:unauthorized"));
    }
    return Promise.reject(error);
  }
);

export const ADMIN_TOKEN_KEY = ADMIN_TOKEN_STORAGE_KEY;
export default adminAxiosClient;
