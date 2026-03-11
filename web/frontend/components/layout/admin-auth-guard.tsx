"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Loader2 } from "lucide-react";
import { useAdminAuthStore } from "@/lib/store/admin-auth-store";

interface AdminAuthGuardProps {
  children: React.ReactNode;
}

export function AdminAuthGuard({ children }: AdminAuthGuardProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading, hydrate, logout } = useAdminAuthStore();

  useEffect(() => {
    hydrate();
  }, [hydrate]);

  useEffect(() => {
    const handleUnauthorized = () => {
      logout();
      router.replace("/admin/login");
    };
    window.addEventListener("admin:unauthorized", handleUnauthorized);
    return () => {
      window.removeEventListener("admin:unauthorized", handleUnauthorized);
    };
  }, [logout, router]);

  useEffect(() => {
    if (isLoading) {
      return;
    }
    if (!isAuthenticated) {
      router.replace("/admin/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="page-background page-background-2 min-h-screen flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
