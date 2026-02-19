"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query-client";

/**
 * Query Provider Component
 * 
 * Wraps the app with TanStack Query's QueryClientProvider.
 * This is a client component (required for React Query).
 * 
 */
export function QueryProvider({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}
