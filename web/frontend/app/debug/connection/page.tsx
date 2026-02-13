'use client';

import { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import axiosClient from '@/lib/axios';

interface ConnectionDiagnostics {
  status: 'idle' | 'loading' | 'success' | 'error';
  statusCode?: number;
  latency?: number;
  headers?: Record<string, string>;
  data?: unknown;
  error?: string;
  errorType?: 'network' | 'cors' | 'http' | 'unknown';
  timestamp?: string;
}

export default function ConnectionDiagnosticsPage() {
  const [diagnostics, setDiagnostics] = useState<ConnectionDiagnostics>({
    status: 'idle',
  });

  const runDiagnostics = useCallback(async () => {
    setDiagnostics({ status: 'loading' });

    const startTime = performance.now();

    try {
      const response = await axiosClient.get('/health');

      const endTime = performance.now();
      const latency = Math.round(endTime - startTime);

      // Extract response headers
      const headers: Record<string, string> = {};
      if (response.headers) {
        Object.keys(response.headers).forEach((key) => {
          headers[key] = String(response.headers[key]);
        });
      }

      setDiagnostics({
        status: 'success',
        statusCode: response.status,
        latency,
        headers,
        data: response.data,
        timestamp: new Date().toISOString(),
      });
    } catch (error: unknown) {
      const endTime = performance.now();
      const latency = Math.round(endTime - startTime);

      let errorMessage = 'Unknown error';
      let errorType: ConnectionDiagnostics['errorType'] = 'unknown';
      let statusCode: number | undefined;

      if (axios.isAxiosError(error)) {
        if (error.code === 'ERR_NETWORK' || error.message.includes('Network Error')) {
          errorMessage = 'Network Error: Unable to connect to the backend API. This usually indicates:\n\n1. Backend server is not running\n2. CORS configuration issue\n3. Wrong API URL (check NEXT_PUBLIC_API_URL)\n4. Network connectivity problem';
          errorType = 'cors';
        } else if (error.response) {
          statusCode = error.response.status;
          errorMessage = `HTTP Error ${statusCode}: ${error.response.statusText}`;
          errorType = 'http';
        } else {
          errorMessage = error.message;
          errorType = 'network';
        }
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }

      setDiagnostics({
        status: 'error',
        statusCode,
        latency,
        error: errorMessage,
        errorType,
        timestamp: new Date().toISOString(),
      });
    }
  }, []);

  useEffect(() => {
    // Auto-run diagnostics on mount in a microtask to satisfy lint rule
    const timer = setTimeout(() => {
      void runDiagnostics();
    }, 0);
    return () => clearTimeout(timer);
  }, [runDiagnostics]);

  const getStatusColor = () => {
    switch (diagnostics.status) {
      case 'success':
        return 'text-green-600 bg-green-50 border-green-200';
      case 'error':
        return 'text-red-600 bg-red-50 border-red-200';
      case 'loading':
        return 'text-blue-600 bg-blue-50 border-blue-200';
      default:
        return 'text-gray-600 bg-gray-50 border-gray-200';
    }
  };

  const getApiUrl = () => {
    return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white shadow-sm rounded-lg p-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Backend Connection Diagnostics
          </h1>
          <p className="text-gray-600 mb-6">
            This page verifies connectivity, CORS configuration, and network communication with the backend API.
          </p>

          {/* API URL Display */}
          <div className="mb-6 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <div className="text-sm font-medium text-gray-700 mb-1">API Base URL</div>
            <div className="text-lg font-mono text-gray-900">{getApiUrl()}</div>
          </div>

          {/* Status Card */}
          <div className={`mb-6 p-4 rounded-lg border ${getStatusColor()}`}>
            <div className="flex items-center justify-between mb-2">
              <div className="text-lg font-semibold">Status</div>
              <button
                onClick={runDiagnostics}
                disabled={diagnostics.status === 'loading'}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
              >
                {diagnostics.status === 'loading' ? 'Testing...' : 'Test Again'}
              </button>
            </div>
            <div className="text-2xl font-bold capitalize">{diagnostics.status}</div>
            {diagnostics.timestamp && (
              <div className="text-sm mt-2 opacity-75">
                Last checked: {new Date(diagnostics.timestamp).toLocaleString()}
              </div>
            )}
          </div>

          {/* Results */}
          {diagnostics.status === 'loading' && (
            <div className="p-6 text-center">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <p className="mt-4 text-gray-600">Testing connection...</p>
            </div>
          )}

          {diagnostics.status === 'success' && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                  <div className="text-sm font-medium text-green-700 mb-1">Status Code</div>
                  <div className="text-2xl font-bold text-green-900">{diagnostics.statusCode}</div>
                </div>
                <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                  <div className="text-sm font-medium text-green-700 mb-1">Latency</div>
                  <div className="text-2xl font-bold text-green-900">{diagnostics.latency}ms</div>
                </div>
              </div>

              <div>
                <h2 className="text-lg font-semibold text-gray-900 mb-2">Response Headers</h2>
                <pre className="p-4 bg-gray-50 rounded-lg border border-gray-200 overflow-x-auto text-sm">
                  {JSON.stringify(diagnostics.headers, null, 2)}
                </pre>
              </div>

              <div>
                <h2 className="text-lg font-semibold text-gray-900 mb-2">Response Data</h2>
                <pre className="p-4 bg-gray-50 rounded-lg border border-gray-200 overflow-x-auto text-sm">
                  {JSON.stringify(diagnostics.data, null, 2)}
                </pre>
              </div>
            </div>
          )}

          {diagnostics.status === 'error' && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                {diagnostics.statusCode && (
                  <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                    <div className="text-sm font-medium text-red-700 mb-1">Status Code</div>
                    <div className="text-2xl font-bold text-red-900">{diagnostics.statusCode}</div>
                  </div>
                )}
                <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                  <div className="text-sm font-medium text-red-700 mb-1">Latency</div>
                  <div className="text-2xl font-bold text-red-900">{diagnostics.latency}ms</div>
                </div>
              </div>

              <div>
                <h2 className="text-lg font-semibold text-red-900 mb-2">Error Details</h2>
                <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                  <div className="text-sm font-medium text-red-700 mb-2">
                    Error Type: {diagnostics.errorType?.toUpperCase()}
                  </div>
                  <pre className="whitespace-pre-wrap text-sm text-red-900">
                    {diagnostics.error}
                  </pre>
                </div>
              </div>

              {diagnostics.errorType === 'cors' && (
                <div className="p-4 bg-yellow-50 rounded-lg border border-yellow-200">
                  <h3 className="font-semibold text-yellow-900 mb-2">CORS Troubleshooting</h3>
                  <ul className="list-disc list-inside space-y-1 text-sm text-yellow-800">
                    <li>Ensure the backend server is running</li>
                    <li>Check backend CORS configuration allows requests from <code className="bg-yellow-100 px-1 rounded">http://localhost:3000</code></li>
                    <li>Verify <code className="bg-yellow-100 px-1 rounded">NEXT_PUBLIC_API_URL</code> is correct</li>
                    <li>Check browser console for additional CORS error details</li>
                  </ul>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
