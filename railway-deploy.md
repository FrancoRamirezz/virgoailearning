# 🚄 Railway Deployment Guide for VirgoLearning Backend

## Quick Deploy to Railway

### 1. **Prepare for Railway**
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login to Railway
railway login

# Initialize project
railway init
```

### 2. **Set Environment Variables**
In Railway dashboard, add these environment variables:

```env
# Database (Railway PostgreSQL)
DATABASE_URL=postgresql://user:pass@host:port/db
DB_HOST=your-railway-postgres-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=railway

# JWT Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-here-64-chars-minimum
JWT_EXPIRY=24h

# Server Configuration
PORT=8080
ENV=production

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com,https://localhost:3000

# Stripe (for payments)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# OAuth (Google)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### 3. **Deploy Commands**
```bash
# Connect to Railway project
railway link

# Add PostgreSQL database
railway add postgresql

# Deploy the backend
railway up
```

## 🎯 What to Demo to Your Team

### **Live Demo Script**

#### **1. Architecture Overview** (2 minutes)
```bash
# Show the production-grade middleware stack
curl -v https://your-app.railway.app/health

# Point out the security headers in response:
# - X-Request-ID: abc123def456
# - X-Content-Type-Options: nosniff
# - X-Frame-Options: DENY
# - Content-Security-Policy: ...
# - Strict-Transport-Security: max-age=31536000
```

#### **2. Rate Limiting Demo** (1 minute)
```bash
# Rapid fire requests to show rate limiting
for i in {1..105}; do curl -s https://your-app.railway.app/health; done
# Should get 429 Too Many Requests after 100 requests
```

#### **3. Structured Logging** (1 minute)
Show Railway logs dashboard:
```json
{
  "timestamp": "2026-01-04T20:39:40Z",
  "request_id": "abc123def456",
  "method": "POST",
  "path": "/api/v1/auth/register",
  "status_code": 201,
  "duration": "150ms",
  "ip": "192.168.1.100",
  "user_agent": "curl/7.68.0"
}
```

#### **4. Validation System** (2 minutes)
```bash
# Show user-friendly validation errors
curl -X POST https://your-app.railway.app/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "weak",
    "first_name": ""
  }'

# Response shows detailed validation:
{
  "success": false,
  "message": "Validation failed", 
  "errors": {
    "email": "email must be a valid email address",
    "password": "password must contain at least 8 characters including uppercase, lowercase, number and special character",
    "first_name": "first_name is required"
  }
}
```

#### **5. Complete User Journey** (3 minutes)
```bash
# 1. Register user
# 2. Login and get JWT
# 3. Create a course (instructor)
# 4. Enroll student
# 5. Process payment
# 6. Show admin endpoints protection
```

## ✅ Production Readiness Checklist

Run these tests locally before deployment:

### **Core Functionality**
```bash
# Run unit tests
TEST_DB_HOST=localhost TEST_DB_USER=sneharoy TEST_DB_PASSWORD="" TEST_DB_NAME=backend_test TEST_JWT_SECRET=test-secret go test ./handlers -v

# Test middleware stack
./test_middleware.sh

# Build check
go build -o /dev/null .
```

### **Security Testing**
```bash
# 1. Security Headers
curl -I https://your-app.railway.app/health | grep -E "X-|Content-Security"

# 2. CORS Testing
curl -H "Origin: https://malicious-site.com" https://your-app.railway.app/api/v1/auth/login

# 3. Rate Limiting
ab -n 150 -c 10 https://your-app.railway.app/health

# 4. SQL Injection Protection (should fail safely)
curl -X POST https://your-app.railway.app/api/v1/auth/login \
  -d '{"email": "admin'\'' OR 1=1--", "password": "anything"}'
```

### **Performance Baseline**
```bash
# Load testing with Apache Bench
ab -n 1000 -c 50 https://your-app.railway.app/health

# Should handle:
# - 1000 requests in <10 seconds
# - No failed requests
# - Average response time <100ms
```

### **Database Testing**
```bash
# Connection pooling
# Migration integrity  
# Transaction handling
go test ./database -v
```

## 🎉 What Makes This Production-Ready

### **Enterprise Features Implemented:**
- ✅ **Panic Recovery** - Server never crashes
- ✅ **Rate Limiting** - Prevents DDoS/abuse
- ✅ **Request Timeouts** - No hanging requests
- ✅ **Security Headers** - OWASP compliance
- ✅ **CORS Protection** - Secure cross-origin requests
- ✅ **Input Validation** - Prevents malicious input
- ✅ **Structured Logging** - Full request traceability
- ✅ **Request Size Limits** - Prevents resource exhaustion

### **Why Railway is Perfect:**
- ✅ **Zero-config PostgreSQL** 
- ✅ **Automatic HTTPS/SSL**
- ✅ **Environment variables**
- ✅ **Real-time logs dashboard**
- ✅ **Metrics and monitoring**
- ✅ **Easy scaling**

## 🚀 Deploy Now!

Your backend is ready for production deployment on Railway. The middleware stack provides enterprise-grade reliability, security, and observability.