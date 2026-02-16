"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";

/**
 * Marketing Navigation
 *
 * Clean navigation for marketing pages with links to login and main sections.
 */
export function MarketingNavigation() {
  return (
    <nav className="portal-surface sticky top-0 z-50 border-x-0 border-t-0">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-9 h-9 rounded-lg flex items-center justify-center shadow-md bg-gradient-to-br from-[var(--gradient-start)] to-[var(--gradient-end)]">
              <span className="text-white font-bold text-sm">35</span>
            </div>
            <span className="font-semibold text-foreground">Entorno 35</span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-2 glass-panel rounded-full px-3 py-1.5 border border-white/70">
            {[
              { href: "#features", label: "Características" },
              { href: "#compliance", label: "Cumplimiento" },
              { href: "#pricing", label: "Precios" },
            ].map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-sm font-medium text-muted-foreground px-3 py-1 rounded-full hover:bg-white/85 hover:text-foreground transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </div>

          {/* CTA Button */}
          <div className="flex items-center space-x-3">
            <Button variant="ghost" className="text-muted-foreground hover:text-foreground" asChild>
              <Link href="/login">Iniciar sesión</Link>
            </Button>
            <Button className="bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] shadow-md hover:brightness-110 hover-shine text-white" asChild>
              <Link href="/login">Comenzar ahora</Link>
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}
