import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  // Security: Disable X-Powered-By header
  poweredByHeader: false,
  
  // Security: Enable strict mode
  reactStrictMode: true,
  
  // Security: Disable source maps in production
  productionBrowserSourceMaps: false,
  
  // Headers are configured in vercel.json for Vercel deployments
  // For local development, headers can be added here if needed
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on'
          },
          {
            key: 'X-Download-Options',
            value: 'noopen'
          },
        ],
      },
    ];
  },
};

export default nextConfig;

