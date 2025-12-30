import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { AuthGuard } from '@/components/layout/auth-guard';

// Mock next/navigation
const mockPush = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock Zustand store
const mockHydrate = vi.fn();
let mockStoreState: {
  isAuthenticated: boolean;
  isLoading: boolean;
  hydrate: typeof mockHydrate;
} = {
  isAuthenticated: false,
  isLoading: false,
  hydrate: mockHydrate,
};

vi.mock('@/lib/store/auth-store', () => ({
  useAuthStore: vi.fn((selector?: unknown) => {
    if (typeof selector === 'function') {
      return selector(mockStoreState);
    }
    return mockStoreState;
  }),
}));

describe('AuthGuard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset to default state
    mockStoreState = {
      isAuthenticated: false,
      isLoading: false,
      hydrate: mockHydrate,
    };
  });

  describe('The Bouncer', () => {
    it('should redirect to /login when user is not authenticated', () => {
      // Mock store state: unauthenticated, not loading
      mockStoreState = {
        isAuthenticated: false,
        isLoading: false,
        hydrate: mockHydrate,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: router.push was called with /login
      expect(mockPush).toHaveBeenCalledWith('/login');
      expect(mockPush).toHaveBeenCalledTimes(1);

      // Assert: Protected Content is NOT in the document
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });
  });

  describe('The VIP', () => {
    it('should render children when user is authenticated', () => {
      // Mock store state: authenticated
      mockStoreState = {
        isAuthenticated: true,
        isLoading: false,
        hydrate: mockHydrate,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: Protected Content IS visible
      expect(screen.getByText('Protected Content')).toBeInTheDocument();

      // Assert: router.push was NOT called
      expect(mockPush).not.toHaveBeenCalled();
    });
  });

  describe('The Waiting Room', () => {
    it('should show loading spinner when isLoading is true', () => {
      // Mock store state: loading
      mockStoreState = {
        isAuthenticated: false,
        isLoading: true,
        hydrate: mockHydrate,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: Loading Spinner is visible
      // Loader2 from lucide-react has specific classes
      const spinner = document.querySelector('.animate-spin');
      expect(spinner).toBeInTheDocument();

      // Assert: Protected Content is NOT in the document
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();

      // Assert: router.push was NOT called (still loading)
      expect(mockPush).not.toHaveBeenCalled();
    });
  });

  describe('Hydration', () => {
    it('should call hydrate on mount', () => {
      mockStoreState = {
        isAuthenticated: true,
        isLoading: false,
        hydrate: mockHydrate,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: hydrate was called on mount
      expect(mockHydrate).toHaveBeenCalledTimes(1);
    });
  });
});
