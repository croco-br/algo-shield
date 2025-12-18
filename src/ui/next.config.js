/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    // Get API URL from environment (for server-side proxying in Docker)
    // In Docker, use the service name 'api', otherwise use localhost
    // This allows Next.js server to proxy API requests to the backend
    const apiUrl = process.env.API_URL || 'http://localhost:8080';
    
    // Only add rewrites if API_URL is set and not empty
    if (!apiUrl || apiUrl.trim() === '') {
      return [];
    }
    
    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
