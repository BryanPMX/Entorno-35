"use client";

import React from "react";
import Image from "next/image";
import { Toaster } from "sonner";
import { QueryProvider } from "@/components/providers/query-provider";

/**
 * Public Assessment Layout
 *
 * Minimalist layout for staff assessment interface.
 * No sidebar, no header, no complex navigation - just the assessment.
 */
export default function PublicAssessmentLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <QueryProvider>
      <div className="page-background page-background-3 relative min-h-screen overflow-hidden">
        <div className="pointer-events-none absolute -left-24 top-8 h-64 w-64 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" />
        <div className="pointer-events-none absolute -right-24 top-16 h-72 w-72 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" />

        {/* Simple Company Logo Header */}
        <div className="portal-surface relative z-10 border-x-0 border-t-0">
          <div className="mx-auto max-w-4xl px-6 py-4">
            <div className="flex items-center justify-center">
              <div className="flex items-center space-x-3">
                <Image src="/logo.png" alt="Entorno 35" width={36} height={36} className="h-9 w-auto object-contain" priority />
                <div>
                  <h1 className="text-lg font-semibold text-foreground">Entorno35</h1>
                  <p className="text-xs text-muted-foreground">Evaluacion NOM-035</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="relative z-10 mx-auto max-w-4xl px-6 py-8">
          {children}
        </div>
      </div>
      <Toaster position="top-center" />
    </QueryProvider>
  );
}
