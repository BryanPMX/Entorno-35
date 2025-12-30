"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query-client";
// Initialize axios client when this client component loads
// This ensures axios is initialized on the client side only
import "@/lib/axios";

/**
 * Query Provider Component
 * 
 * Wraps the app with TanStack Query's QueryClientProvider.
 * This is a client component (required for React Query).
 * 
 * Also initializes the axios client (which logs initialization to console).
 */
export function QueryProvider({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}
