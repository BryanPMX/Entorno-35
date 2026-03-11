"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { LogOut, ClipboardList } from "lucide-react";
import { AdminAuthGuard } from "@/components/layout/admin-auth-guard";
import { useAdminAuthStore } from "@/lib/store/admin-auth-store";
import { Button } from "@/components/ui/button";

export default function AdminProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { logout, email } = useAdminAuthStore();

  const handleLogout = () => {
    logout();
    router.push("/admin/login");
  };

  return (
    <AdminAuthGuard>
      <div className="page-background page-background-2 portal-grid relative min-h-screen">
        <header className="portal-surface sticky top-0 z-20 border-x-0 border-t-0">
          <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6">
            <div className="flex items-center gap-4">
              <p className="text-sm font-semibold text-foreground">Billing Admin Console</p>
              <Link
                href="/admin/refunds"
                className="inline-flex items-center gap-2 rounded-md border border-input bg-background/70 px-3 py-1.5 text-xs font-medium text-foreground hover:bg-accent"
              >
                <ClipboardList className="h-3.5 w-3.5" />
                Reembolsos
              </Link>
            </div>

            <div className="flex items-center gap-3">
              <p className="hidden text-xs text-muted-foreground sm:block">{email ?? "admin"}</p>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleLogout}
                className="gap-2 border-primary/20 bg-background/70"
              >
                <LogOut className="h-4 w-4" />
                Salir
              </Button>
            </div>
          </div>
        </header>

        <main className="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6">{children}</main>
      </div>
    </AdminAuthGuard>
  );
}
