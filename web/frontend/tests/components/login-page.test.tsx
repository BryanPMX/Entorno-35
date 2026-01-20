import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LoginPage from '@/app/(auth)/login/page';
import type { AuthResponse } from '@/types/backend';

// Mock next/navigation
const mockPush = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock auth service - define mock function inside factory for hoisting
vi.mock('@/services/auth.service', () => {
  const mockLogin = vi.fn();
  return {
    authService: {
      login: mockLogin,
    },
  };
});

// Import after mock to get access to mocked module
import { authService } from '@/services/auth.service';

// Mock Zustand store
const mockStoreLogin = vi.fn();
let mockStoreState: {
  login: typeof mockStoreLogin;
} = {
  login: mockStoreLogin,
};

vi.mock('@/lib/store/auth-store', () => ({
  useAuthStore: vi.fn(() => mockStoreState),
}));

// Mock toast (sonner)
vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockStoreState = {
      login: mockStoreLogin,
    };
  });

  describe('Validation', () => {
    it('should show validation errors when submitting empty form', async () => {
      const user = userEvent.setup();
      render(<LoginPage />);

      // Find and click the Sign In button without filling form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Assert: Validation errors appear
      await waitFor(() => {
        expect(screen.getByText(/rfc is required/i)).toBeInTheDocument();
      });

      // Assert: authService.login was NOT called
      expect(vi.mocked(authService.login)).not.toHaveBeenCalled();
    });
  });

  describe('Success Flow', () => {
    it('should successfully login with COMPANY type and redirect to dashboard', async () => {
      const user = userEvent.setup();
      const mockToken = 'mock-jwt-token-123';
      const mockResponse: AuthResponse = { token: mockToken };

      // Mock successful login
      vi.mocked(authService.login).mockResolvedValue(mockResponse);

      render(<LoginPage />);

      // Fill form (Company is default)
      const identifierInput = screen.getByLabelText(/rfc/i);
      await user.type(identifierInput, 'ABC123456789');

      // Submit form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Wait for async operations
      await waitFor(() => {
        // Assert: authService.login was called with correct arguments
        expect(vi.mocked(authService.login)).toHaveBeenCalledTimes(1);
        expect(vi.mocked(authService.login)).toHaveBeenCalledWith({
          identifier: 'ABC123456789',
          type: 'COMPANY',
        });
      });

      // Assert: store.login was called
      expect(mockStoreLogin).toHaveBeenCalledTimes(1);
      expect(mockStoreLogin).toHaveBeenCalledWith(mockResponse);

      // Assert: router.push('/dashboard') was called
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/dashboard');
      });
    });

  });


  describe('Error Handling', () => {
    it('should show error toast when login fails', async () => {
      const user = userEvent.setup();
      const { toast } = await import('sonner');

      // Mock failed login
      vi.mocked(authService.login).mockRejectedValue(new Error('Invalid credentials'));

      render(<LoginPage />);

      // Fill and submit form
      const identifierInput = screen.getByLabelText(/rfc/i);
      await user.type(identifierInput, 'INVALID123');

      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Wait for error handling
      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith(
          'Invalid credentials. Please try again.'
        );
      });

      // Assert: router.push was NOT called
      expect(mockPush).not.toHaveBeenCalled();
    });
  });
});
