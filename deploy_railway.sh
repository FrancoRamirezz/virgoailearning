#!/bin/bash

echo "🚀 VirgoLearning Backend - Railway Deployment"
echo "============================================="

# Check if logged in
echo "📋 Checking Railway login status..."
if ! railway whoami; then
    echo "❌ Please login to Railway first:"
    echo "   railway login"
    exit 1
fi

echo "✅ Railway CLI authenticated"

# Check if project is linked
echo "📋 Checking project status..."
if ! railway status; then
    echo "❌ No Railway project linked. Please run:"
    echo "   railway init"
    echo "   railway add postgresql"
    exit 1
fi

# Build test
echo "📋 Testing build..."
if ! go build -o main .; then
    echo "❌ Build failed. Please fix build errors first."
    exit 1
fi
rm -f main
echo "✅ Build test passed"

# Check for required env vars
echo "📋 Environment Variables Setup:"
echo "⚠️  MANUAL STEP REQUIRED:"
echo "   1. Go to Railway dashboard"
echo "   2. Add these environment variables:"
echo "      - JWT_SECRET: your-super-secure-jwt-secret-key-32-chars-minimum"
echo "      - JWT_EXPIRY: 24h"
echo "      - ENV: production"  
echo "      - PORT: 8080"
echo "      - CORS_ALLOWED_ORIGINS: https://your-frontend-domain.com"
echo ""
echo "🔍 Railway will auto-configure database variables (DATABASE_URL, etc.)"
echo ""

read -p "Have you added the environment variables? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Please add environment variables first, then run this script again"
    exit 1
fi

# Deploy
echo "🚀 Deploying to Railway..."
railway up

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 DEPLOYMENT SUCCESSFUL!"
    echo "========================="
    echo ""
    echo "📊 Next Steps:"
    echo "1. Get your Railway URL: railway open"
    echo "2. Test health endpoint: curl https://your-app.railway.app/health"
    echo "3. Test API endpoints: curl https://your-app.railway.app/api/v1/courses"
    echo "4. Check logs: railway logs"
    echo ""
    echo "🎯 Your VirgoLearning Backend is LIVE!"
else
    echo "❌ Deployment failed. Check logs with: railway logs"
    exit 1
fi