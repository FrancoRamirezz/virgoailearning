/** @type {import('next').NextConfig} */
const nextConfig = {
/* this  */
 source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*'
};

export default nextConfig;
