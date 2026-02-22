"use client";

import Link from "next/link";
import Image from "next/image";
import { Button } from "@/components/ui/button";

/**
 * Marketing Navigation
 *
 * Clean navigation for marketing pages with links to login and main sections.
 */
export function MarketingNavigation() {
  return (
    <nav className="portal-surface sticky top-0 z-50 border-x-0 border-t-0">
      <div className="w-full px-4 sm:px-6 lg:px-10">
        <div className="grid h-16 grid-cols-[1fr_auto_1fr] items-center gap-3">
          {/* Logo */}
          <Link href="/" className="justify-self-start">
            <Image src="/logo.png" alt="Entorno 35" width={812} height={293} className="h-11 w-auto object-contain md:h-14" priority />
          </Link>

          {/* Navigation Links */}
          <div className="hidden justify-self-center md:flex md:items-center md:space-x-2 portal-surface rounded-full px-3 py-1.5 border-0">
            {[
              { href: "#features", label: "Características" },
              { href: "#compliance", label: "Cumplimiento" },
              { href: "#pricing", label: "Precios" },
            ].map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="portal-nav-link text-sm font-medium px-3 py-1 rounded-full transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </div>

          {/* Auth Buttons */}
          <div className="flex items-center gap-2 justify-self-end">
            <Button variant="outline" className="hidden sm:inline-flex" asChild>
              <Link href="/register">Crear cuenta</Link>
            </Button>
            <Button variant="ghost" className="text-muted-foreground hover:text-foreground" asChild>
              <Link href="/login">Iniciar sesión</Link>
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}
