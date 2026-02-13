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
    <nav className="sticky top-0 z-50 border-b border-white/40 bg-white/70 backdrop-blur-2xl shadow-[0_8px_24px_rgba(15,23,42,0.08)] dark:border-white/10 dark:bg-background/70">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-9 h-9 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center shadow-md">
              <span className="text-white font-bold text-sm">35</span>
            </div>
            <span className="font-semibold text-gray-900 dark:text-white">Entorno 35</span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-2 glass-panel rounded-full px-3 py-1.5 border border-white/70 dark:border-white/10">
            {[
              { href: "#features", label: "Características" },
              { href: "#compliance", label: "Cumplimiento" },
              { href: "#pricing", label: "Precios" },
            ].map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-sm font-medium text-gray-700 dark:text-gray-200 px-3 py-1 rounded-full hover:bg-white/85 hover:text-gray-900 dark:hover:bg-white/10 transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </div>

          {/* CTA Button */}
          <div className="flex items-center space-x-3">
            <Button variant="ghost" className="text-gray-700 dark:text-gray-200" asChild>
              <Link href="/login">Iniciar sesión</Link>
            </Button>
            <Button className="bg-gradient-to-r from-blue-600 to-purple-600 shadow-md hover:brightness-110 hover-shine" asChild>
              <Link href="/login">Comenzar ahora</Link>
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}
