"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/lib/store/auth-store";
import { AuthGuard } from "@/components/layout/auth-guard";
import { Button } from "@/components/ui/button";
import { LogOut, LayoutDashboard, Users, FileText } from "lucide-react";
import { toast } from "sonner";

/**
 * Dashboard Layout
 * 
 * Protected layout for authenticated users with:
 * - AuthGuard wrapper (redirects to /login if not authenticated)
 * - Sidebar navigation (left, fixed)
 * - Header with logout button (top)
 */
export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    toast.success("Logged out successfully");
    router.push("/login");
  };

  return (
    <AuthGuard>
      <div className="min-h-screen flex flex-col bg-gray-50">
        {/* Header */}
        <header className="bg-white border-b border-gray-200 sticky top-0 z-10">
          <div className="flex items-center justify-between px-6 h-16">
            <h1 className="text-xl font-semibold text-gray-900">
              Entorno 35
            </h1>
            <Button
              variant="outline"
              size="sm"
              onClick={handleLogout}
              className="gap-2"
            >
              <LogOut className="h-4 w-4" />
              Logout
            </Button>
          </div>
        </header>

        <div className="flex flex-1">
          {/* Sidebar */}
          <aside className="w-64 bg-white border-r border-gray-200 fixed h-[calc(100vh-4rem)] top-16 overflow-y-auto">
            <nav className="p-4 space-y-2">
              <Link
                href="/dashboard"
                className="flex items-center gap-3 px-4 py-3 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <LayoutDashboard className="h-5 w-5" />
                Dashboard
              </Link>
              <Link
                href="/dashboard/staff"
                className="flex items-center gap-3 px-4 py-3 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <Users className="h-5 w-5" />
                Staff
              </Link>
              <Link
                href="/dashboard/assessments"
                className="flex items-center gap-3 px-4 py-3 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <FileText className="h-5 w-5" />
                Assessments
              </Link>
            </nav>
          </aside>

          {/* Main Content */}
          <main className="flex-1 ml-64 p-6">
            {children}
          </main>
        </div>
      </div>
    </AuthGuard>
  );
}

