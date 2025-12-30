import axiosClient from "@/lib/axios";
import type { AuthResponse, LoginRequest } from "@/types/backend";

/**
 * Authentication Service
 * 
 * Service layer for authentication operations.
 * Components should NEVER call axios directly - they call these service methods.
 * This centralizes error handling and type validation.
 */
class AuthService {
  /**
   * Login
   * Authenticates a user (company or staff) and returns a JWT token
   * 
   * @param credentials - Login credentials (identifier, type, optional company_id)
   * @returns Promise resolving to AuthResponse with token
   * @throws AxiosError on authentication failure
   */
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await axiosClient.post<AuthResponse>("/auth/login", credentials);
    
    // Store token in localStorage
    if (response.data.token && typeof window !== "undefined") {
      localStorage.setItem("token", response.data.token);
    }
    
    return response.data;
  }

  /**
   * Logout
   * Clears the stored authentication token
   */
  logout(): void {
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
  }

  /**
   * Get Token
   * Retrieves the stored authentication token
   * 
   * @returns Token string or null if not found
   */
  getToken(): string | null {
    if (typeof window !== "undefined") {
      return localStorage.getItem("token");
    }
    return null;
  }

  /**
   * Check if user is authenticated
   * 
   * @returns true if token exists, false otherwise
   */
  isAuthenticated(): boolean {
    return this.getToken() !== null;
  }
}

// Export singleton instance
export const authService = new AuthService();
export default authService;

