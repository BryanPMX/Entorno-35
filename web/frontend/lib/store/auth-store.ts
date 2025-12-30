import { create } from "zustand";
import type { User, AuthResponse } from "@/types/backend";
import { decodeJWT } from "@/lib/jwt";

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (data: AuthResponse) => void;
  logout: () => void;
  hydrate: () => void;
}

/**
 * Global Auth Store (Zustand)
 * 
 * Manages authentication state across the application.
 * - Stores user info and token
 * - Handles login/logout
 * - Hydrates state from localStorage on mount
 */
export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: true,

  /**
   * Login action
   * Updates state with token and decoded user info
   */
  login: (data: AuthResponse) => {
    const token = data.token;
    
    // Decode user from JWT token
    const user = decodeJWT(token);
    
    // Save to localStorage
    if (typeof window !== "undefined") {
      localStorage.setItem("token", token);
    }

    set({
      token,
      user,
      isAuthenticated: true,
      isLoading: false,
    });
  },

  /**
   * Logout action
   * Clears state and removes token from localStorage
   */
  logout: () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }

    set({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
    });
  },

  /**
   * Hydrate action
   * Checks localStorage on mount and restores auth state
   */
  hydrate: () => {
    if (typeof window === "undefined") {
      set({ isLoading: false });
      return;
    }

    const token = localStorage.getItem("token");
    
    if (token) {
      // Decode user from token
      const user = decodeJWT(token);
      
      set({
        token,
        user,
        isAuthenticated: true,
        isLoading: false,
      });
    } else {
      set({
        token: null,
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });
    }
  },
}));

