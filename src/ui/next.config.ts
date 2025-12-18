import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    // Get API URL from environment (for server-side proxying in Docker)
    // In Docker, use the service name 'api', otherwise use localhost
    const apiUrl = process.env.API_URL || 'http://localhost:8080';
    
    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
