import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@/lib/store/auth-store";

/**
 * API Base URL
 * In production, this should come from environment variables
 */
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/**
 * Singleton Axios Instance
 * 
 * This instance handles:
 * - Automatic token injection from localStorage
 * - Global error handling (401 redirects)
 * - Base URL configuration
 */
const axiosClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * Request Interceptor
 * Automatically injects Authorization header with token from localStorage
 */
axiosClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Check for token in localStorage (client-side only)
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("token");
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

/**
 * Response Interceptor
 * Handles global error responses, especially 401 Unauthorized
 */
axiosClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // Handle 401 Unauthorized - redirect to login
    if (error.response?.status === 401) {
      // Only redirect if not already on login page and not making a login request
      if (typeof window !== "undefined") {
        const isLoginRequest = error.config?.url?.includes("/auth/login");
        const isLoginPage = window.location.pathname === "/login";

        if (!isLoginRequest) {
          useAuthStore.getState().logout();

          if (!isLoginPage) {
            window.dispatchEvent(new Event("auth:unauthorized"));
          }
        }
      }
    }

    return Promise.reject(error);
  }
);

// Log initialization (for verification)
if (typeof window !== "undefined") {
  console.log("✅ Axios client initialized with base URL:", API_BASE_URL);
}

export default axiosClient;
