import { create } from "zustand";
import { jwtDecode, type JwtPayload } from "jwt-decode";
import { ADMIN_TOKEN_KEY } from "@/lib/admin-axios";

interface AdminJWTPayload extends JwtPayload {
  role?: string;
  email?: string;
}

interface AdminAuthState {
  token: string | null;
  email: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string) => void;
  logout: () => void;
  hydrate: () => void;
}

function parseAdminToken(token: string): { email: string | null; exp: number | null; isAdmin: boolean } | null {
  try {
    const payload = jwtDecode<AdminJWTPayload>(token);
    const role = typeof payload.role === "string" ? payload.role.toLowerCase() : "";
    return {
      email: typeof payload.email === "string" ? payload.email : null,
      exp: typeof payload.exp === "number" ? payload.exp : null,
      isAdmin: role === "admin",
    };
  } catch {
    return null;
  }
}

function isExpired(exp: number | null): boolean {
  if (exp === null) {
    return true;
  }
  const now = Math.floor(Date.now() / 1000);
  return exp <= now;
}

export const useAdminAuthStore = create<AdminAuthState>((set) => ({
  token: null,
  email: null,
  isAuthenticated: false,
  isLoading: true,

  login: (token: string) => {
    const parsed = parseAdminToken(token);
    if (!parsed || !parsed.isAdmin || isExpired(parsed.exp)) {
      if (typeof window !== "undefined") {
        localStorage.removeItem(ADMIN_TOKEN_KEY);
      }
      set({
        token: null,
        email: null,
        isAuthenticated: false,
        isLoading: false,
      });
      return;
    }

    if (typeof window !== "undefined") {
      localStorage.setItem(ADMIN_TOKEN_KEY, token);
    }

    set({
      token,
      email: parsed.email,
      isAuthenticated: true,
      isLoading: false,
    });
  },

  logout: () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem(ADMIN_TOKEN_KEY);
    }
    set({
      token: null,
      email: null,
      isAuthenticated: false,
      isLoading: false,
    });
  },

  hydrate: () => {
    if (typeof window === "undefined") {
      set({ isLoading: false });
      return;
    }

    const token = localStorage.getItem(ADMIN_TOKEN_KEY);
    if (!token) {
      set({
        token: null,
        email: null,
        isAuthenticated: false,
        isLoading: false,
      });
      return;
    }

    const parsed = parseAdminToken(token);
    if (!parsed || !parsed.isAdmin || isExpired(parsed.exp)) {
      localStorage.removeItem(ADMIN_TOKEN_KEY);
      set({
        token: null,
        email: null,
        isAuthenticated: false,
        isLoading: false,
      });
      return;
    }

    set({
      token,
      email: parsed.email,
      isAuthenticated: true,
      isLoading: false,
    });
  },
}));
