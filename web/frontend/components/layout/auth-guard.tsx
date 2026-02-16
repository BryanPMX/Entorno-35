"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Loader2 } from "lucide-react";
import { useAuthStore } from "@/lib/store/auth-store";
import { isTokenExpired } from "@/lib/jwt";

interface AuthGuardProps {
  children: React.ReactNode;
}

/**
 * AuthGuard Component
 * 
 * Protects routes by ensuring the user is authenticated.
 * - Calls hydrate() on mount to check localStorage
 * - Shows loading spinner while checking authentication
 * - Redirects to /login if not authenticated
 * - Renders children if authenticated
 */
export function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading, token, hydrate, logout } = useAuthStore();
  const tokenExpired = token ? isTokenExpired(token) : true;

  useEffect(() => {
    // Hydrate auth state from localStorage on mount
    hydrate();
  }, [hydrate]);

  useEffect(() => {
    const handleUnauthorized = () => {
      logout();
      router.replace("/login");
    };

    window.addEventListener("auth:unauthorized", handleUnauthorized);
    return () => {
      window.removeEventListener("auth:unauthorized", handleUnauthorized);
    };
  }, [logout, router]);

  useEffect(() => {
    // Redirect to login when token is missing/invalid/expired (after hydrate)
    if (isLoading) {
      return;
    }

    if (!isAuthenticated || !token) {
      router.replace("/login");
      return;
    }

    if (tokenExpired) {
      logout();
      router.replace("/login");
    }
  }, [isAuthenticated, isLoading, logout, router, token, tokenExpired]);

  // Show loading spinner while checking authentication
  if (isLoading) {
    return (
      <div className="portal-shell min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  // Don't render anything while redirecting
  if (!isAuthenticated || !token || tokenExpired) {
    return null;
  }

  // Render protected content
  return <>{children}</>;
}
