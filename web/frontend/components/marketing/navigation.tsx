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
    <nav className="border-b bg-white/95 backdrop-blur-sm sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-blue-700 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">35</span>
            </div>
            <span className="font-semibold text-gray-900">Entorno 35</span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-8">
            <Link
              href="#features"
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              Características
            </Link>
            <Link
              href="#compliance"
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              Cumplimiento
            </Link>
            <Link
              href="#pricing"
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              Precios
            </Link>
          </div>

          {/* CTA Button */}
          <div className="flex items-center space-x-4">
            <Button variant="ghost" asChild>
              <Link href="/login">Iniciar Sesión</Link>
            </Button>
            <Button asChild>
              <Link href="/login">Comenzar Ahora</Link>
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}