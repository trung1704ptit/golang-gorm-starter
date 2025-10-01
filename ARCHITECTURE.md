# Travel Blog Backend - Architecture Documentation

## Overview
This project follows a clean architecture pattern with clear separation of concerns between layers.

## Architecture Layers

### 📁 Project Structure
```
travel-blog-backend2/
├── controllers/          # HTTP handlers (thin layer)
├── services/            # Business logic layer
├── models/              # Data models & DTOs
├── routes/              # Route definitions
├── middleware/          # HTTP middleware
├── initializers/        # App initialization
├── utils/               # Helper functions
└── main.go             # Application entry point
```

## Layer Responsibilities

### 1️⃣ Models Layer (`models/`)
**Purpose:** Define data structures and DTOs

- `user.model.go` - User entity and related DTOs (SignUpInput, SignInInput, UserResponse)
- `post.model.go` - Post entity and related DTOs (CreatePostRequest, UpdatePost)

**Responsibilities:**
- Database entity definitions
- Request/Response DTOs
- Data validation tags

---

### 2️⃣ Services Layer (`services/`) 🆕
**Purpose:** Contains all business logic

#### `auth.service.go`
- `SignUp(payload)` - User registration logic
- `SignIn(payload)` - User authentication
- `RefreshAccessToken(token)` - Token refresh logic
- `GetUserByID(id)` - User retrieval

#### `user.service.go`
- `GetUserResponse(user)` - User data transformation
- Future: UpdateUserProfile, ChangePassword, etc.

#### `post.service.go`
- `CreatePost(payload, userID)` - Post creation
- `UpdatePost(postID, payload, userID)` - Post updates
- `FindPostByID(postID)` - Single post retrieval
- `FindPosts(page, limit)` - Paginated posts
- `DeletePost(postID)` - Post deletion

**Responsibilities:**
- Business logic implementation
- Database operations via GORM
- Data validation at business level
- Error handling with meaningful messages
- Reusable across different interfaces (REST, GraphQL, gRPC, etc.)

---

### 3️⃣ Controllers Layer (`controllers/`)
**Purpose:** Handle HTTP requests and responses (thin layer)

#### `auth.controller.go`
- Binds HTTP request JSON to DTOs
- Calls AuthService methods
- Maps service errors to HTTP status codes
- Sets cookies for authentication
- Returns HTTP responses

#### `user.controller.go`
- Extracts current user from context
- Calls UserService methods
- Returns user data as HTTP response

#### `post.controller.go`
- Handles HTTP request parsing
- Calls PostService methods
- Maps errors to appropriate HTTP responses
- Handles pagination parameters

**Responsibilities:**
- HTTP request/response handling
- Request validation (JSON binding)
- Error-to-HTTP status mapping
- Cookie management
- NO business logic

---

### 4️⃣ Routes Layer (`routes/`)
**Purpose:** Define API endpoints and middleware

- `auth.routes.go` - Authentication endpoints
- `user.routes.go` - User endpoints
- `post.routes.go` - Post endpoints

**Responsibilities:**
- Route definitions
- Middleware attachment
- Controller method mapping

---

### 5️⃣ Middleware Layer (`middleware/`)
**Purpose:** Request preprocessing

- `deserialize-user.go` - JWT authentication middleware

**Responsibilities:**
- Authentication/Authorization
- Request preprocessing
- Logging, rate limiting, etc.

---

## Data Flow

```
HTTP Request
    ↓
Routes (route definition)
    ↓
Middleware (auth, validation)
    ↓
Controller (HTTP handling)
    ↓
Service (business logic)
    ↓
Database (via GORM)
    ↓
Service (transform data)
    ↓
Controller (HTTP response)
    ↓
HTTP Response
```

## Dependency Injection Flow

```go
main.go:
  DB Connection (initializers.DB)
      ↓
  Services (new with DB)
      ↓
  Controllers (new with Services)
      ↓
  Route Controllers (new with Controllers)
      ↓
  Gin Router (register routes)
```

## Benefits of This Architecture

### ✅ Separation of Concerns
- Each layer has a single, well-defined responsibility
- Controllers only handle HTTP, Services handle business logic

### ✅ Testability
- Services can be unit tested independently
- Controllers can be tested with mocked services
- No need for HTTP server to test business logic

### ✅ Reusability
- Services can be used by different interfaces (REST, GraphQL, CLI)
- Business logic is centralized and not duplicated

### ✅ Maintainability
- Easy to locate and fix bugs (clear layer boundaries)
- Easy to add new features (extend services)
- Easy to understand codebase structure

### ✅ Scalability
- Can easily add caching layer
- Can add repository pattern between services and DB
- Can split services into microservices if needed

## Example Usage

### Adding a New Feature

**Example: Add "Update User Profile" feature**

1. **Add DTO** in `models/user.model.go`:
```go
type UpdateProfileInput struct {
    Name  string `json:"name" binding:"required"`
    Photo string `json:"photo"`
}
```

2. **Add Service Method** in `services/user.service.go`:
```go
func (s *UserService) UpdateProfile(userID uuid.UUID, payload *models.UpdateProfileInput) (*models.UserResponse, error) {
    // Business logic here
    var user models.User
    s.DB.First(&user, "id = ?", userID)
    
    user.Name = payload.Name
    user.Photo = payload.Photo
    s.DB.Save(&user)
    
    return s.GetUserResponse(&user), nil
}
```

3. **Add Controller Method** in `controllers/user.controller.go`:
```go
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
    currentUser := ctx.MustGet("currentUser").(models.User)
    var payload *models.UpdateProfileInput
    
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userResponse, err := uc.userService.UpdateProfile(currentUser.ID, payload)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"data": userResponse})
}
```

4. **Add Route** in `routes/user.routes.go`:
```go
router.PUT("/profile", middleware.DeserializeUser(), uc.userController.UpdateProfile)
```

## Testing Strategy

### Unit Tests (Services)
```go
func TestAuthService_SignUp(t *testing.T) {
    // Mock DB
    // Test business logic
    // Assert results
}
```

### Integration Tests (Controllers)
```go
func TestAuthController_SignUp(t *testing.T) {
    // Mock service
    // Test HTTP handling
    // Assert HTTP responses
}
```

### E2E Tests
- Test complete flow with real database
- Use test containers for PostgreSQL

## Future Improvements

1. **Repository Pattern**: Add repository layer between services and DB
2. **Dependency Injection Container**: Use wire or fx
3. **Caching Layer**: Add Redis for caching
4. **Event-Driven**: Add event publishing for async operations
5. **API Versioning**: Add v1, v2 routes
6. **GraphQL Support**: Reuse services for GraphQL resolvers

## Conclusion

This architecture provides a solid foundation for building scalable, maintainable, and testable applications. The clear separation of concerns makes it easy to understand, extend, and maintain the codebase.

---

**Author:** Travel Blog Backend Team  
**Last Updated:** October 2025  
**Version:** 2.0

