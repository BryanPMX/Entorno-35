"use client";

import React from "react";

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
    <div className="min-h-screen bg-slate-50">
      {/* Simple Company Logo Header */}
      <div className="bg-white border-b border-slate-200 shadow-sm">
        <div className="max-w-2xl mx-auto px-6 py-4">
          <div className="flex items-center justify-center">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                <span className="text-primary-foreground font-bold text-sm">35</span>
              </div>
              <div>
                <h1 className="text-lg font-semibold text-slate-900">Entorno35</h1>
                <p className="text-xs text-slate-500">NOM-035 Assessment</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-2xl mx-auto px-6 py-8">
        {children}
      </div>
    </div>
  );
}