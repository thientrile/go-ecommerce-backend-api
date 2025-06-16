# API Documentation

This document provides comprehensive information about the API endpoints, authentication, and usage examples.

## Base URLs

| Environment | Base URL |
|-------------|----------|
| Local Development | `http://crm.shopdev.test/v1/2025` |
| Docker Environment | `http://crm.shopdev.com/v1/2025` |
| Production | `https://your-domain.com/v1/2025` |

## Authentication

### JWT Authentication
The API uses JSON Web Tokens (JWT) for authentication.

#### Login Process
1. Send login credentials to `/auth/login`
2. Receive JWT token in response
3. Include token in subsequent requests

```bash
# Login request
curl -X POST http://crm.shopdev.test/v1/2025/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response
{
  "code": 200,
  "message": "Success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

#### Using JWT Token
Include the token in the Authorization header:

```bash
curl -X GET http://crm.shopdev.test/v1/2025/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## API Endpoints

### Authentication Endpoints

#### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "confirm_password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Registration successful",
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "status": "pending_verification"
  }
}
```

#### POST /auth/login
Login with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "jwt_token_here",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

#### POST /auth/verify
Verify email with OTP code.

**Request Body:**
```json
{
  "email": "user@example.com",
  "verify_code": "123456",
  "verify_type": "email"
}
```

#### POST /auth/logout
Logout current user session.

**Headers:**
```
Authorization: Bearer jwt_token_here
```

### User Management Endpoints

#### GET /users/profile
Get current user profile.

**Headers:**
```
Authorization: Bearer jwt_token_here
```

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "avatar": "https://example.com/avatar.jpg",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### PUT /users/profile
Update user profile.

**Headers:**
```
Authorization: Bearer jwt_token_here
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "phone": "+1234567890"
}
```

#### POST /users/change-password
Change user password.

**Headers:**
```
Authorization: Bearer jwt_token_here
```

**Request Body:**
```json
{
  "old_password": "old_password123",
  "new_password": "new_password123",
  "confirm_password": "new_password123"
}
```

### Health Check Endpoints

#### GET /health
Check application health status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T00:00:00Z",
  "services": {
    "database": "connected",
    "redis": "connected",
    "kafka": "connected"
  },
  "version": "1.0.0"
}
```

## Response Format

### Success Response
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Bad Request",
  "error": "Detailed error message",
  "details": {
    "field": "validation error details"
  }
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 422 | Validation Error |
| 500 | Internal Server Error |

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Guest users**: 100 requests per hour
- **Authenticated users**: 1000 requests per hour
- **Admin users**: 10000 requests per hour

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1609459200
```

## Validation Rules

### Email Validation
- Must be a valid email format
- Must be unique in the system
- Maximum length: 255 characters

### Password Validation
- Minimum length: 8 characters
- Must contain at least one uppercase letter
- Must contain at least one lowercase letter
- Must contain at least one number
- Must contain at least one special character

### Phone Validation
- Must be a valid phone number format
- International format preferred (+country_code)

## Example Usage

### Complete Registration Flow

```bash
# 1. Register new user
curl -X POST http://crm.shopdev.test/v1/2025/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "SecurePass123!",
    "confirm_password": "SecurePass123!",
    "first_name": "Jane",
    "last_name": "Doe"
  }'

# 2. Verify email (check email for OTP)
curl -X POST http://crm.shopdev.test/v1/2025/auth/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "verify_code": "123456",
    "verify_type": "email"
  }'

# 3. Login
curl -X POST http://crm.shopdev.test/v1/2025/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "SecurePass123!"
  }'

# 4. Use JWT token for authenticated requests
curl -X GET http://crm.shopdev.test/v1/2025/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Error Handling Example

```javascript
// JavaScript example with error handling
async function apiRequest(endpoint, options = {}) {
  try {
    const response = await fetch(`http://crm.shopdev.test/v1/2025${endpoint}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        ...options.headers
      },
      ...options
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'API request failed');
    }

    return data;
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
}

// Usage
try {
  const profile = await apiRequest('/users/profile');
  console.log('User profile:', profile.data);
} catch (error) {
  console.error('Failed to fetch profile:', error.message);
}
```

## Swagger Documentation

Interactive API documentation is available via Swagger UI:

- **Local**: http://crm.shopdev.test/swagger/index.html
- **Docker**: http://crm.shopdev.com/swagger/index.html

The Swagger documentation includes:
- All available endpoints
- Request/response schemas
- Interactive testing interface
- Authentication requirements
- Example requests and responses

## SDK and Client Libraries

### Go Client
```go
package main

import (
    "fmt"
    "github.com/your-org/go-ecommerce-client"
)

func main() {
    client := ecommerce.NewClient("http://crm.shopdev.test/v1/2025")
    
    // Login
    token, err := client.Auth.Login("user@example.com", "password123")
    if err != nil {
        panic(err)
    }
    
    // Set authentication token
    client.SetToken(token)
    
    // Get user profile
    profile, err := client.Users.GetProfile()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Welcome, %s!\n", profile.FirstName)
}
```

### JavaScript Client
```javascript
class EcommerceClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.token = localStorage.getItem('auth_token');
  }

  async login(email, password) {
    const response = await fetch(`${this.baseURL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    const data = await response.json();
    if (data.code === 200) {
      this.token = data.data.token;
      localStorage.setItem('auth_token', this.token);
    }
    return data;
  }

  async getProfile() {
    const response = await fetch(`${this.baseURL}/users/profile`, {
      headers: { 'Authorization': `Bearer ${this.token}` }
    });
    return response.json();
  }
}

// Usage
const client = new EcommerceClient('http://crm.shopdev.test/v1/2025');
```

## Testing

### Postman Collection
Import the Postman collection for easy API testing:
1. Download: `docs/postman/ecommerce-api.json`
2. Import into Postman
3. Set environment variables
4. Run test cases

### curl Examples
All API endpoints can be tested using curl commands provided in this documentation.

### Automated Testing
Run the test suite:
```bash
# Run all API tests
make test-api

# Run specific test category
go test ./tests/api/auth_test.go
go test ./tests/api/user_test.go
```
