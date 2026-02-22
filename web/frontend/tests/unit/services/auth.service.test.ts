import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import type { LoginRequest, AuthResponse } from '@/types/backend';

// Mock axios client - must use factory function for hoisting
vi.mock('@/lib/axios', () => {
  const mockPost = vi.fn();
  return {
    default: {
      post: mockPost,
    },
  };
});

// Import after mock
import { authService } from '@/services/auth.service';
import axiosClient from '@/lib/axios';

function toBase64Url(value: string): string {
  return btoa(value).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
}

function createToken(payload: Record<string, unknown>): string {
  const header = toBase64Url(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const body = toBase64Url(JSON.stringify(payload));
  return `${header}.${body}.signature`;
}

describe('AuthService', () => {
  // Mock localStorage
  const localStorageMock = (() => {
    let store: Record<string, string> = {};

    return {
      getItem: (key: string) => store[key] || null,
      setItem: (key: string, value: string) => {
        store[key] = value.toString();
      },
      removeItem: (key: string) => {
        delete store[key];
      },
      clear: () => {
        store = {};
      },
    };
  })();

  beforeEach(() => {
    // Setup localStorage mock
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock,
      writable: true,
    });
    localStorageMock.clear();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('login', () => {
    it('should call axios.post with correct endpoint and payload', async () => {
      const mockResponse: AuthResponse = {
        token: 'test-token-123',
      };

      const mockCredentials: LoginRequest = {
        identifier: 'ABC123456789',
        type: 'COMPANY',
        password: 'secure-password-123',
      };

      // Mock successful response
      vi.mocked(axiosClient.post).mockResolvedValue({
        data: mockResponse,
      } as unknown as Awaited<ReturnType<typeof axiosClient.post>>);

      const result = await authService.login(mockCredentials);

      // Verify axios.post was called with correct arguments
      expect(axiosClient.post).toHaveBeenCalledTimes(1);
      expect(axiosClient.post).toHaveBeenCalledWith('/auth/login', mockCredentials);

      // Verify result
      expect(result).toEqual(mockResponse);
    });

    it('should save token to localStorage after successful login', async () => {
      const mockToken = 'test-token-456';
      const mockResponse: AuthResponse = {
        token: mockToken,
      };

      const mockCredentials: LoginRequest = {
        identifier: 'ABC123456789',
        type: 'COMPANY',
        password: 'secure-password-123',
      };

      vi.mocked(axiosClient.post).mockResolvedValue({
        data: mockResponse,
      } as unknown as Awaited<ReturnType<typeof axiosClient.post>>);

      await authService.login(mockCredentials);

      // Verify token was saved to localStorage
      expect(localStorage.getItem('token')).toBe(mockToken);
    });
  });

  describe('logout', () => {
    it('should remove token from localStorage', () => {
      // Set a token first
      localStorage.setItem('token', 'test-token');

      authService.logout();

      expect(localStorage.getItem('token')).toBeNull();
    });
  });

  describe('getToken', () => {
    it('should return token from localStorage if exists', () => {
      const token = 'test-token-789';
      localStorage.setItem('token', token);

      expect(authService.getToken()).toBe(token);
    });

    it('should return null if token does not exist', () => {
      localStorage.removeItem('token');

      expect(authService.getToken()).toBeNull();
    });
  });

  describe('isAuthenticated', () => {
    it('should return true if token exists and is not expired', () => {
      const validToken = createToken({
        company_id: 'company-1',
        role: 'company',
        exp: Math.floor(Date.now() / 1000) + 3600,
      });
      localStorage.setItem('token', validToken);

      expect(authService.isAuthenticated()).toBe(true);
    });

    it('should return false if token does not exist', () => {
      localStorage.removeItem('token');

      expect(authService.isAuthenticated()).toBe(false);
    });

    it('should return false if token is expired', () => {
      const expiredToken = createToken({
        company_id: 'company-1',
        role: 'company',
        exp: Math.floor(Date.now() / 1000) - 10,
      });
      localStorage.setItem('token', expiredToken);

      expect(authService.isAuthenticated()).toBe(false);
    });
  });
});
