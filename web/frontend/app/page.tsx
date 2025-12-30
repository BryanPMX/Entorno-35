"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/store/auth-store";

/**
 * Root Page
 * 
 * Redirects to /dashboard if authenticated, otherwise to /login
 */
export default function HomePage() {
  const router = useRouter();
  const { isAuthenticated, hydrate } = useAuthStore();

  useEffect(() => {
    // Hydrate auth state
    hydrate();
  }, [hydrate]);

  useEffect(() => {
    // Redirect based on authentication status
    if (!isAuthenticated) {
      router.push("/login");
    } else {
      router.push("/dashboard");
    }
  }, [isAuthenticated, router]);

  // Show nothing while redirecting
  return null;
}
