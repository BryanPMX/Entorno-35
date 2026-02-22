"use client";

import "@/styles/tw-animate.css";
import Link from "next/link";
import Image from "next/image";
import { usePathname, useRouter } from "next/navigation";
import { LogOut, LayoutDashboard, Users, FileText } from "lucide-react";
import { Toaster, toast } from "sonner";
import { AuthGuard } from "@/components/layout/auth-guard";
import { QueryProvider } from "@/components/providers/query-provider";
import { Button } from "@/components/ui/button";
import { translations } from "@/lib/translations";
import { useAuthStore } from "@/lib/store/auth-store";
import { cn } from "@/lib/utils";

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
  const pathname = usePathname();
  const router = useRouter();
  const { logout } = useAuthStore();

  const navItems = [
    {
      href: "/dashboard",
      label: translations.nav.dashboard,
      icon: LayoutDashboard,
      isActive: pathname === "/dashboard",
    },
    {
      href: "/dashboard/staff",
      label: translations.nav.staff,
      icon: Users,
      isActive: pathname.startsWith("/dashboard/staff"),
    },
    {
      href: "/dashboard/assessments",
      label: translations.nav.assessments,
      icon: FileText,
      isActive: pathname.startsWith("/dashboard/assessments"),
    },
  ];

  const handleLogout = () => {
    logout();
    toast.success(translations.nav.logoutSuccess);
    router.push("/login");
  };

  return (
    <QueryProvider>
      <AuthGuard>
        <div className="page-background page-background-2 portal-grid relative flex min-h-screen flex-col">
          <div className="pointer-events-none absolute -left-24 top-20 h-72 w-72 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" />
          <div className="pointer-events-none absolute -right-28 top-10 h-80 w-80 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" />

          {/* Header */}
          <header className="portal-surface sticky top-0 z-20 border-x-0 border-t-0">
            <div className="flex items-center justify-between px-6 h-16">
              <div className="flex items-center">
                <Image src="/logo.png" alt="Entorno 35" width={812} height={293} className="h-12 w-auto object-contain md:h-14" priority />
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="gap-2 border-primary/20 bg-background/70 text-foreground hover:bg-accent/70"
              >
                <LogOut className="h-4 w-4" />
                {translations.nav.logout}
              </Button>
            </div>
          </header>

          <div className="flex flex-1">
            {/* Sidebar */}
            <aside className="portal-surface fixed top-16 h-[calc(100vh-4rem)] w-64 overflow-y-auto border-y-0 border-l-0">
              <nav className="p-4 space-y-2">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  return (
                    <Link
                      key={item.href}
                      href={item.href}
                      className={cn(
                        "portal-nav-link flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-all",
                        item.isActive && "portal-nav-link-active"
                      )}
                    >
                      <Icon className="h-5 w-5" />
                      {item.label}
                    </Link>
                  );
                })}
              </nav>
            </aside>

            {/* Main Content */}
            <main className="relative z-10 ml-64 flex-1 p-6">
              {children}
            </main>
          </div>
        </div>
      </AuthGuard>
      <Toaster position="top-center" />
    </QueryProvider>
  );
}
