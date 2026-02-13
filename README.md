# Event Booking REST API

A full-featured event booking REST API built with Go. I created this project to deepen my understanding of the Go programming language and backend development best practices.

## 🎯 Learning Goals

This project helped me learn and practice:
- Building RESTful APIs in Go
- Clean architecture (handlers → services → database)
- JWT authentication and authorization
- Database operations with SQLite
- Middleware implementation
- Error handling and validation
- Environment configuration
- Security best practices (password hashing, token management)

## ✨ Features

- **User Management**
  - User registration with password hashing
  - Login with JWT token generation
  - User authentication middleware

- **Event Management**
  - Create, read, update, and delete events
  - Authorization checks (users can only modify their own events)
  - Event listing

- **Event Registration**
  - Users can register for events
  - Cancel event registrations
  - Track registered users per event

## 🛠 Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Gin
- **Database**: SQLite3
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator
- **Configuration**: godotenv
- **Testing**: testify, go.uber.org/mock

## 📁 Project Structure

```
event-booking/
├── main.go                 # Application entry point
├── .env                    # Environment variables (not committed)
├── db/
│   ├── db.go              # Database initialization
│   ├── events.go          # Event database operations
│   ├── events_test.go     # Event repository tests
│   ├── users.go           # User database operations
│   ├── users_test.go      # User repository tests
│   ├── register.go        # Registration database operations
│   ├── register_test.go   # Registration repository tests
│   └── testdb.go          # Test database helpers
├── models/
│   ├── event.go           # Event model
│   ├── user.go            # User model
│   └── register.go        # Registration model
├── routes/
│   ├── routes.go          # Route registration
│   ├── events.go          # Event handlers
│   ├── users.go           # User handlers
│   └── register.go        # Registration handlers
├── services/
│   ├── event.go           # Event business logic
│   ├── event_test.go      # Event service tests
│   ├── user.go            # User business logic
│   ├── user_test.go       # User service tests
│   ├── register.go        # Registration business logic
│   ├── register_test.go   # Registration service tests
│   └── mocks/             # Generated mock repositories
├── middleware/
│   └── auth.go            # JWT authentication middleware
├── utils/
│   ├── hash.go            # Password hashing utilities
│   ├── hash_test.go       # Password hashing tests
│   ├── jwt.go             # JWT token utilities
│   ├── jwt_test.go        # JWT token tests
│   ├── validators.go      # Custom validation functions
│   └── validators_test.go # Validation tests
└── testutil/
    └── env.go             # Test environment setup
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or higher
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd event-bookng
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
Create a `.env` file in the root directory:
```env
PORT=8000
JWT_SECRET=your-super-secret-key-change-this
```

4. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8000`

## 🧪 Testing

This project includes comprehensive unit and integration tests with **78 tests** achieving over **90% coverage** of core business logic.

### Test Structure

```
✅ utils:     92.3% coverage (24 tests)
   - Password hashing and validation
   - JWT token generation and verification
   - Custom date validators

✅ services:  98.2% coverage (33 tests)
   - EventService: CRUD + authorization checks
   - UserService: User management + JWT login
   - EventRegisterService: Event registration workflows

✅ db:        80.6% coverage (21 tests)
   - EventRepository: Full CRUD operations
   - UserRepository: User CRUD + password hashing
   - RegisterRepository: Registration management
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests for specific package
go test ./services/... -v
go test ./db/... -v
go test ./utils/... -v

# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Patterns

- **Service Layer**: Uses `go.uber.org/mock` to mock repository interfaces
- **Repository Layer**: Uses in-memory SQLite for integration testing
- **Utils Layer**: Pure function tests with no dependencies

All tests follow the Arrange-Act-Assert pattern and are fully isolated.

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/signup` | Register a new user | No |
| POST | `/login` | Login and get JWT token | No |

### Events

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/events` | Get all events | Yes |
| GET | `/events/:id` | Get event by ID | Yes |
| POST | `/events` | Create a new event | Yes |
| PUT | `/events/:id` | Update an event | Yes (owner only) |
| DELETE | `/events/:id` | Delete an event | Yes (owner only) |

### Event Registration

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/events/:id/register` | Register for an event | Yes |
| DELETE | `/events/:id/register` | Cancel event registration | Yes |

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users` | Get all users | Yes |

## ✅ Validation Rules

All requests are automatically validated. Invalid data returns `400 Bad Request` with error details.

### User Validation
- **Email**: Must be a valid email format
- **Password**: Minimum 8 characters

### Event Validation
- **Name**: Required, 3-100 characters
- **Description**: Required, 5-500 characters
- **Location**: Required, 3-100 characters
- **DateTime**: Required, must be a future date/time

## 🔐 Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Get your token by logging in via the `/login` endpoint.

## 📝 Example Requests

### Register a User
```http
POST http://localhost:8000/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepass123"
}
```

**Note**: Password must be at least 8 characters.

### Login
```http
POST http://localhost:8000/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Create an Event
```http
POST http://localhost:8000/events
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "name": "Go Conference 2026",
  "description": "Annual Go programming conference with workshops and talks",
  "location": "San Francisco, CA",
  "dateTime": "2026-06-15T09:00:00Z"
}
```

**Note**: All fields are required. The date must be in the future, and field lengths must meet validation requirements.

## 🏗 Architecture

This project follows a clean, layered architecture with **dependency injection** for loose coupling:

```
Handlers → Services → Repositories → Database
```

### Layers

1. **Handlers (Routes)**: Handle HTTP requests/responses, parse input, return JSON
2. **Services**: Contain business logic and authorization rules
3. **Repositories**: Abstract database operations
4. **Models**: Define data structures
5. **Middleware**: Handle cross-cutting concerns (authentication)
6. **Utils**: Provide reusable utilities (JWT, password hashing)

### Dependency Injection

The application uses **dependency injection** to achieve loose coupling between layers:

- **Handlers** receive **Services** as dependencies
- **Services** receive **Repositories** as dependencies
- **Repositories** use the database connection

**Benefits:**
- **Testability**: Each layer can be tested in isolation with mocked dependencies
- **Flexibility**: Easy to swap implementations (e.g., SQLite → PostgreSQL)
- **Maintainability**: Clear separation of concerns and explicit dependencies

## 🔒 Security Features

- Password hashing with bcrypt
- JWT-based authentication
- Authorization checks for resource ownership
- Environment-based configuration
- Comprehensive input validation (email format, field lengths, custom validators)
- SQL injection prevention (parameterized queries)

## 📖 What I Learned

Through building this project, I gained hands-on experience with:

- **Go Fundamentals**: Structs, interfaces, error handling, packages
- **HTTP/REST**: Building RESTful APIs with proper status codes and responses
- **Database**: Working with SQL databases in Go
- **Security**: Implementing authentication, authorization, and data protection
- **Validation**: Request validation with struct tags and custom validators
- **Dependency Injection**: Implementing DI for loose coupling and testability
- **Repository Pattern**: Abstracting database operations for flexibility
- **Architecture**: Organizing code for maintainability and scalability
- **Testing**: Unit tests, integration tests, mocking, and test-driven development
- **Test Patterns**: Arrange-Act-Assert, table-driven tests, and test isolation
- **Best Practices**: Clean code, separation of concerns, and error handling patterns

## 🤝 Contributing

This is a learning project, but suggestions and feedback are welcome!

## 📄 License

This project is open source and available for learning purposes.

---

**Built with ❤️ while learning Go**
