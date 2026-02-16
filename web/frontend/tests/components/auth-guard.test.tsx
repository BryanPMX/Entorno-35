import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { AuthGuard } from '@/components/layout/auth-guard';

function toBase64Url(value: string): string {
  return btoa(value).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
}

function createToken(payload: Record<string, unknown>): string {
  const header = toBase64Url(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const body = toBase64Url(JSON.stringify(payload));
  return `${header}.${body}.signature`;
}

// Mock next/navigation
const mockReplace = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    replace: mockReplace,
  }),
}));

// Mock Zustand store
const mockHydrate = vi.fn();
const mockLogout = vi.fn();
let mockStoreState: {
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
  hydrate: typeof mockHydrate;
  logout: typeof mockLogout;
} = {
  isAuthenticated: false,
  isLoading: false,
  token: null,
  hydrate: mockHydrate,
  logout: mockLogout,
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
      token: null,
      hydrate: mockHydrate,
      logout: mockLogout,
    };
  });

  describe('The Bouncer', () => {
    it('should redirect to /login when user is not authenticated', () => {
      // Mock store state: unauthenticated, not loading
      mockStoreState = {
        isAuthenticated: false,
        isLoading: false,
        token: null,
        hydrate: mockHydrate,
        logout: mockLogout,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: router.replace was called with /login
      expect(mockReplace).toHaveBeenCalledWith('/login');
      expect(mockReplace).toHaveBeenCalledTimes(1);

      // Assert: Protected Content is NOT in the document
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });

    it('should logout and redirect to /login when token is expired', () => {
      mockStoreState = {
        isAuthenticated: true,
        isLoading: false,
        token: createToken({
          company_id: 'company-1',
          role: 'company',
          exp: Math.floor(Date.now() / 1000) - 10,
        }),
        hydrate: mockHydrate,
        logout: mockLogout,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      expect(mockLogout).toHaveBeenCalledTimes(1);
      expect(mockReplace).toHaveBeenCalledWith('/login');
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });
  });

  describe('The VIP', () => {
    it('should render children when user is authenticated', () => {
      // Mock store state: authenticated
      mockStoreState = {
        isAuthenticated: true,
        isLoading: false,
        token: createToken({
          company_id: 'company-1',
          role: 'company',
          exp: Math.floor(Date.now() / 1000) + 3600,
        }),
        hydrate: mockHydrate,
        logout: mockLogout,
      };

      render(
        <AuthGuard>
          <div>Protected Content</div>
        </AuthGuard>
      );

      // Assert: Protected Content IS visible
      expect(screen.getByText('Protected Content')).toBeInTheDocument();

      // Assert: router.replace was NOT called
      expect(mockReplace).not.toHaveBeenCalled();
    });
  });

  describe('The Waiting Room', () => {
    it('should show loading spinner when isLoading is true', () => {
      // Mock store state: loading
      mockStoreState = {
        isAuthenticated: false,
        isLoading: true,
        token: null,
        hydrate: mockHydrate,
        logout: mockLogout,
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

      // Assert: router.replace was NOT called (still loading)
      expect(mockReplace).not.toHaveBeenCalled();
    });
  });

  describe('Hydration', () => {
    it('should call hydrate on mount', () => {
      mockStoreState = {
        isAuthenticated: true,
        isLoading: false,
        token: createToken({
          company_id: 'company-1',
          role: 'company',
          exp: Math.floor(Date.now() / 1000) + 3600,
        }),
        hydrate: mockHydrate,
        logout: mockLogout,
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
