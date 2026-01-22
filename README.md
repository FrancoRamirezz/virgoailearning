# Backend Template

A comprehensive Go backend template with authentication, database integration, and RESTful API endpoints.

## Features

- **Authentication**: JWT-based authentication with login/register
- **Database**: PostgreSQL with GORM ORM
- **API**: RESTful API with proper error handling
- **Middleware**: CORS, logging, and authentication middleware
- **Configuration**: Environment-based configuration
- **Security**: Password hashing, JWT tokens, role-based access

## Project Structure

```
backend/
├── config/          # Configuration management
├── database/        # Database connection and migrations
├── handlers/        # HTTP request handlers
├── middleware/      # Custom middleware
├── models/          # Data models
├── routes/          # Route definitions
├── utils/           # Utility functions
├── .env.example     # Environment variables template
├── go.mod          # Go module file
├── go.sum          # Go dependencies
└── main.go         # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.25.1 or later
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Set up PostgreSQL database:
```sql
CREATE DATABASE backend_db;
```

5. Run the application:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your .env file).

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get current user profile (requires auth)

### Users
- `GET /api/v1/users` - Get all users (requires auth)
- `GET /api/v1/users/{id}` - Get user by ID (requires auth)
- `PUT /api/v1/users/{id}` - Update user (requires auth)
- `DELETE /api/v1/users/{id}` - Delete user (requires auth)

### Posts
- `GET /api/v1/posts` - Get all posts (requires auth)
- `GET /api/v1/posts/{id}` - Get post by ID (requires auth)
- `POST /api/v1/posts` - Create new post (requires auth)
- `PUT /api/v1/posts/{id}` - Update post (requires auth, author or admin)
- `DELETE /api/v1/posts/{id}` - Delete post (requires auth, author or admin)

### Health Check
- `GET /health` - Server health check

## Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | password |
| `DB_NAME` | Database name | backend_db |
| `JWT_SECRET` | JWT secret key | your-super-secret-jwt-key-here |
| `JWT_EXPIRY` | JWT expiry duration | 24h |
| `PORT` | Server port | 8080 |
| `ENV` | Environment | development |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | http://localhost:3000 |

## Database Models

### User
- `id` (Primary Key)
- `email` (Unique)
- `password` (Hashed)
- `first_name`
- `last_name`
- `role` (user/admin)
- `is_active`
- `created_at`
- `updated_at`
- `deleted_at` (Soft delete)

### Post
- `id` (Primary Key)
- `title`
- `content`
- `author_id` (Foreign Key to User)
- `is_published`
- `created_at`
- `updated_at`
- `deleted_at` (Soft delete)

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o backend main.go
```

### Docker (Optional)
```bash
# Build image
docker build -t backend .

# Run container
docker run -p 8080:8080 backend
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.


