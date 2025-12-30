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
        expect(screen.getByText(/identifier is required/i)).toBeInTheDocument();
      });

      // Assert: authService.login was NOT called
      expect(vi.mocked(authService.login)).not.toHaveBeenCalled();
    });

    it('should require company_id for STAFF type', async () => {
      const user = userEvent.setup();
      render(<LoginPage />);

      // Switch to STAFF tab
      const staffButton = screen.getByRole('button', { name: /staff/i });
      await user.click(staffButton);

      // Fill identifier but not company_id
      const identifierInput = screen.getByLabelText(/curp/i);
      await user.type(identifierInput, 'TESTCURP12345678901234');

      // Submit form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Assert: Validation error for company_id appears
      await waitFor(() => {
        expect(
          screen.getByText(/company id is required for staff login/i)
        ).toBeInTheDocument();
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

    it('should successfully login with STAFF type and redirect to dashboard', async () => {
      const user = userEvent.setup();
      const mockToken = 'mock-jwt-token-456';
      const mockResponse: AuthResponse = { token: mockToken };

      // Mock successful login
      vi.mocked(authService.login).mockResolvedValue(mockResponse);

      render(<LoginPage />);

      // Switch to STAFF tab
      const staffButton = screen.getByRole('button', { name: /staff/i });
      await user.click(staffButton);

      // Fill form
      const identifierInput = screen.getByLabelText(/curp/i);
      await user.type(identifierInput, 'TESTCURP12345678901234');

      const companyIdInput = screen.getByLabelText(/company id/i);
      await user.type(companyIdInput, '550e8400-e29b-41d4-a716-446655440000');

      // Submit form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Wait for async operations
      await waitFor(() => {
        // Assert: authService.login was called with correct arguments
        expect(vi.mocked(authService.login)).toHaveBeenCalledTimes(1);
        expect(vi.mocked(authService.login)).toHaveBeenCalledWith({
          identifier: 'TESTCURP12345678901234',
          type: 'STAFF',
          company_id: '550e8400-e29b-41d4-a716-446655440000',
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

  describe('Tab Switching', () => {
    it('should change input label when switching to STAFF tab', async () => {
      const user = userEvent.setup();
      render(<LoginPage />);

      // Initially, should show RFC label (COMPANY is default)
      expect(screen.getByLabelText(/rfc/i)).toBeInTheDocument();

      // Click Staff button
      const staffButton = screen.getByRole('button', { name: /staff/i });
      await user.click(staffButton);

      // Assert: Label changes to CURP
      expect(screen.getByLabelText(/curp/i)).toBeInTheDocument();
      expect(screen.queryByLabelText(/rfc/i)).not.toBeInTheDocument();

      // Assert: Company ID field appears
      expect(screen.getByLabelText(/company id/i)).toBeInTheDocument();
    });

    it('should send STAFF type in payload when STAFF tab is selected', async () => {
      const user = userEvent.setup();
      const mockToken = 'mock-jwt-token-789';
      const mockResponse: AuthResponse = { token: mockToken };

      vi.mocked(authService.login).mockResolvedValue(mockResponse);

      render(<LoginPage />);

      // Switch to STAFF tab
      const staffButton = screen.getByRole('button', { name: /staff/i });
      await user.click(staffButton);

      // Fill form
      const identifierInput = screen.getByLabelText(/curp/i);
      await user.type(identifierInput, 'TESTCURP12345678901234');

      const companyIdInput = screen.getByLabelText(/company id/i);
      await user.type(companyIdInput, '550e8400-e29b-41d4-a716-446655440000');

      // Submit form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Wait for async operations
      await waitFor(() => {
        // Assert: Payload sent to authService has type: "STAFF"
        expect(vi.mocked(authService.login)).toHaveBeenCalledWith(
          expect.objectContaining({
            type: 'STAFF',
          })
        );
      });
    });

    it('should send COMPANY type in payload when COMPANY tab is selected', async () => {
      const user = userEvent.setup();
      const mockToken = 'mock-jwt-token-012';
      const mockResponse: AuthResponse = { token: mockToken };

      vi.mocked(authService.login).mockResolvedValue(mockResponse);

      render(<LoginPage />);

      // COMPANY is default, but let's explicitly click it
      const companyButton = screen.getByRole('button', { name: /company/i });
      await user.click(companyButton);

      // Fill form
      const identifierInput = screen.getByLabelText(/rfc/i);
      await user.type(identifierInput, 'ABC123456789');

      // Submit form
      const signInButton = screen.getByRole('button', { name: /sign in/i });
      await user.click(signInButton);

      // Wait for async operations
      await waitFor(() => {
        // Assert: Payload sent to authService has type: "COMPANY"
        expect(vi.mocked(authService.login)).toHaveBeenCalledWith(
          expect.objectContaining({
            type: 'COMPANY',
          })
        );
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
