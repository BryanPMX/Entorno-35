import { QueryClient } from "@tanstack/react-query";

/**
 * TanStack Query Client Configuration
 * 
 * Configured to prevent race conditions and reduce unnecessary refetches:
 * - refetchOnWindowFocus: false - Prevents refetching when user switches tabs
 * - refetchOnMount: true - Ensures fresh data on component mount
 * - retry: 1 - Retry failed requests once (reduces network noise)
 * - staleTime: 5 minutes - Data is considered fresh for 5 minutes
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false, // Prevents refetch on tab focus (reduces noise)
      refetchOnMount: true, // Always refetch on mount for fresh data
      retry: 1, // Retry once on failure
      staleTime: 1000 * 60 * 5, // 5 minutes - data is fresh for this duration
      gcTime: 1000 * 60 * 10, // 10 minutes - cache garbage collection time (formerly cacheTime)
    },
    mutations: {
      retry: 1, // Retry once on mutation failure
    },
  },
});

