
# Blog Platform Documentation

## Overview
Blog Platform is a RESTful API built with Go (Gin) and PostgreSQL. It supports user registration, authentication, blog publishing, tagging, and commenting. JWT is used for authentication, and email verification is required for new users.

---

## Project Structure
```
blog-platform/
├── delivery/         # HTTP layer: main entry, controllers, routers
├── domain/           # Core business models and interfaces
├── infrastructure/   # Services: JWT, email, password, middleware
├── repositories/     # Data access layer
├── usecases/         # Business logic
├── test/             # Unit and integration tests
├── docs/             # Documentation
```

---

## Main Features
- User registration, login, profile view/update, and email activation
- Blog creation, listing, and detail view with tags
- JWT authentication and secure password hashing
- Role-based and owner-based route protection

---

## API Endpoints & Examples

### User
- `POST /register` — Register a new user
	- Request: `{ "username": "johndoe", "email": "john@example.com", "password": "Password123!" }`
- `POST /login` — Login and receive tokens
	- Request: `{ "identifier": "johndoe", "password": "Password123!" }`
	- Response: `{ "access": "<JWT>", "refresh": "<JWT>", "message": "Logged in successfully" }`
- `GET /users/:id` — Get user profile (auth required)
- `PATCH /users/:id` — Update user profile (auth required, owner only)

### Blog
- `POST /blogs` — Create a new blog (auth required)
	- Request: `{ "title": "My First Blog", "content": "...", "tags": "go,api,blog" }`
- `GET /blogs` — List all blogs
- `GET /blogs/:id` — Get a blog by ID

---

## Environment Variables
Set these in a `.env` file:
```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=blogdb
DB_PORT=5432
JWT_ACCESS_SECRET=youraccesssecret
JWT_REFRESH_SECRET=yourrefreshsecret
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=yourpassword
SMTP_FROM=your@email.com
PROTOCOL=http
DOMAIN=localhost
PORT=8080
```

---

## Authentication Flow
1. Register and receive activation email
2. Activate account via email link
3. Login to receive JWT tokens
4. Use `Authorization: Bearer <token>` for protected routes

---

## Error Handling
All errors are returned as JSON with an `error` field and appropriate HTTP status code.
Example: `{ "error": "invalid credentials" }`

---

## Running & Testing
1. Set up PostgreSQL and configure `.env`
2. Run: `go run delivery/main.go`
3. API available at configured port
4. Run all tests: `go test ./...`

---

## Extending the Platform
- Add features by extending domain interfaces and implementing new usecases, repositories, and controllers
- Add new routes in the `routers` package
