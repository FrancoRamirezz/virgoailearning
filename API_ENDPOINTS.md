# 🚀 VirgoLearning Backend - Complete API Documentation

## 🌐 Base URL
- **Local**: `http://localhost:8080`
- **Production**: `https://your-app.railway.app`

## 📋 Complete Endpoint List

### **Health Check**
- `GET /health` - Server health status

### **🔐 Authentication Routes** (`/api/v1/auth`)
**Public Routes:**
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user  
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password with token
- `GET /auth/google` - Initiate Google OAuth
- `GET /auth/google/callback` - Google OAuth callback

**Protected Routes (require JWT):**
- `GET /auth/profile` - Get current user profile
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/change-password` - Change password
- `POST /auth/logout` - Logout user

### **👥 User Management** (`/api/v1/users`) - Protected
- `GET /users` - List all users
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### **📚 Course Management** (`/api/v1/courses`)
**Public Routes:**
- `GET /courses` - List published courses
- `GET /courses/{id}` - Get course details
- `GET /courses/{id}/content` - Get course content (for enrolled users)

**Protected Routes (require instructor/admin):**
- `POST /courses` - Create new course
- `PUT /courses/{id}` - Update course
- `DELETE /courses/{id}` - Delete course
- `POST /courses/{id}/publish` - Publish course
- `POST /courses/{id}/content` - Add course content
- `PUT /courses/{id}/content/{contentId}` - Update course content
- `DELETE /courses/{id}/content/{contentId}` - Delete course content
- `POST /courses/{id}/enroll` - Enroll in course

### **📝 Posts/Blog** (`/api/v1/posts`) - Protected
- `GET /posts` - List all posts
- `GET /posts/{id}` - Get post by ID
- `POST /posts` - Create new post
- `PUT /posts/{id}` - Update post
- `DELETE /posts/{id}` - Delete post

### **📞 Contact Management** (`/api/v1/contacts`)
**Public Routes:**
- `POST /contacts` - Submit contact form

**Protected Routes:**
- `GET /contacts` - List contacts
- `GET /contacts/{id}` - Get contact by ID
- `GET /contacts/search` - Search contacts
- `PUT /contacts/{id}` - Update contact
- `DELETE /contacts/{id}` - Delete contact

### **💳 Payment System** (`/api/v1/payments`) - Protected

**Course Payments:**
- `POST /payments/process` - Process course payment
- `GET /payments/history` - Get payment history
- `GET /payments/{id}` - Get payment details
- `GET /payments/{id}/status` - Check payment status
- `GET /payments/{id}/receipt` - Get payment receipt

**Stripe Integration:**
- `POST /payments/intent` - Create payment intent
- `GET /payments/intent/{id}` - Get payment intent
- `POST /payments/confirm` - Confirm payment
- `POST /payments/webhook` - Stripe webhook (public)

**Admin Payment Routes:**
- `GET /admin/payments` - Get all payments
- `POST /admin/payments/{id}/refund` - Refund payment
- `GET /admin/payments/stats` - Payment statistics

## 🔒 Authentication & Authorization

### **JWT Token Usage**
Include in header for protected routes:
```
Authorization: Bearer <your-jwt-token>
```

### **User Roles**
- **`user`** - Regular students (default)
- **`instructor`** - Can create/manage courses  
- **`admin`** - Full system access

### **Route Protection**
- **Public Routes** - No authentication required
- **Protected Routes** - Require valid JWT token
- **Instructor Routes** - Require instructor or admin role
- **Admin Routes** - Require admin role

## 📊 Example Requests

### **Register User**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### **Login**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "SecurePass123!"
  }'
```

### **Create Course (Instructor)**
```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Programming",
    "description": "Learn programming basics",
    "price": 9999,
    "currency": "usd",
    "category": "programming",
    "level": "beginner",
    "duration": 40
  }'
```

### **Enroll in Course**
```bash
curl -X POST http://localhost:8080/api/v1/courses/1/enroll \
  -H "Authorization: Bearer <token>"
```

### **Process Payment**
```bash
curl -X POST http://localhost:8080/api/v1/payments/process \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": 1,
    "amount": 9999,
    "currency": "usd"
  }'
```

## 🛡️ Security Features

### **Middleware Protection**
- ✅ **Request ID Tracking** - Every request gets unique ID
- ✅ **Panic Recovery** - Server never crashes
- ✅ **Security Headers** - OWASP compliance
- ✅ **Rate Limiting** - 100 requests/minute per IP
- ✅ **Request Timeout** - 30-second limit
- ✅ **CORS Protection** - Secure cross-origin requests
- ✅ **Input Validation** - User-friendly error messages
- ✅ **Request Size Limits** - 10MB maximum

### **Validation Examples**
Invalid registration returns detailed errors:
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": "email must be a valid email address",
    "password": "password must contain at least 8 characters including uppercase, lowercase, number and special character"
  }
}
```

## 🚀 Ready for Production!

**Complete API with:**
- ✅ **48+ endpoints** across all features
- ✅ **Role-based authorization** 
- ✅ **Input validation** on all inputs
- ✅ **Security middleware** protection
- ✅ **Payment processing** with Stripe
- ✅ **Course management** system
- ✅ **User authentication** (JWT + OAuth)
- ✅ **Contact form** handling
- ✅ **Admin panel** access

**Test with:** `./test_routes.sh`