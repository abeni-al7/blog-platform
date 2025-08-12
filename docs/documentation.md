# Blog Platform Backend API Documentation

## Quickstart

Follow these steps to get the API running locally.

1) Prerequisites
- Go 1.20+ and PostgreSQL running locally or reachable

2) Create .env (project root) - follow .env.example

3) Run the server (auto-migration enabled)
```sh
go run delivery/main.go
```

4) Use the API
- Base URL: http://localhost:8000
- Get tokens: POST /register, then POST /login
- Send Authorization header: `Authorization: Bearer <ACCESS_TOKEN>` for protected routes

Postman collection: https://documenter.getpostman.com/view/42847133/2sB3BEoViy

## Project Overview

The Blog Platform backend is a robust RESTful API designed to power a modern blogging application. It provides secure user authentication, blog management, tagging, commenting, and administrative features. The backend is built for extensibility, maintainability, and security, enabling developers to build, customize, and scale the platform as needed.

### Main Features
- User registration, login, profile management, and email activation
- JWT-based authentication and role-based authorization (admin, user)
- Blog CRUD (create, read, update, delete) with tagging
- Pagination, filtering, and full-text search for blogs
- Views, likes/unlikes, and popularity metrics per blog
- Commenting system for blogs (create + list with pagination)
- Password reset and account recovery
- Admin controls for user promotion/demotion
- AI-powered blog idea generation and improvement suggestions (optional)

---
## You can access the Postman documentation here:[Postman docs](https://documenter.getpostman.com/view/42847133/2sB3BEoViy)
## Architecture

- **Style:** Layered architecture, RESTful API
- **Layers:**
  - **Delivery:** HTTP layer (Gin controllers, routers)
  - **Domain:** Core business models and interfaces
  - **Usecases:** Business logic and application rules
  - **Repositories:** Data access and persistence
  - **Infrastructure:** Services (JWT, email, password, middleware)

### Key Components
- **Controllers:** Handle HTTP requests and responses
- **Routers:** Define API routes and apply middleware
- **Domain Models:** User, Blog, Tag, Comment, Token
- **Usecases:** Implement business logic for users, blogs, etc.
- **Repositories:** Interact with the database using GORM
- **Middleware:** Auth, admin, account owner, error handling
- **External Services:** SMTP email, AI service (optional)

---

## Technology Stack

- **Language:** Go
- **Framework:** Gin (HTTP server)
- **ORM:** GORM
- **Database:** PostgreSQL
- **Email:** SMTP
- **Authentication:** JWT (github.com/golang-jwt/jwt)
- **Password Hashing:** bcrypt

---

## API Endpoints

Base URL (local): http://localhost:8080

Common headers for protected routes:
- Authorization: Bearer <ACCESS_TOKEN>
- Content-Type: application/json

### Users

| Method | URL                        | Auth         | Description                       |
|--------|----------------------------|--------------|-----------------------------------|
| POST   | /register                  | No           | Register a new user               |
| POST   | /login                     | No           | Login and receive tokens          |
| POST   | /token/refresh             | Yes           | Refresh JWT tokens                |
| POST   | /logout                    | Yes          | Logout (invalidate tokens)        |
| POST   | /reset-password            | Yes          | Change password (logged in)       |
| POST   | /forgot-password           | No           | Request password reset email      |
| POST   | /password/:id/update       | Yes           | Set new password (via token)      |
| GET    | /users/:id                 | Owner/Admin  | Get user profile                  |
| PATCH  | /users/:id                 | Owner/Admin  | Update user profile               |
| PUT    | /users/:id/promote         | Admin        | Promote user to admin             |
| PUT    | /users/:id/demote          | Admin        | Demote user from admin            |

#### Example: Register
Request:
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "Password123!"
}
```
Response:
```json
{
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "status": "inactive"
  }
}
```

#### Example: Login
Request:
```json
{
  "identifier": "johndoe",
  "password": "Password123!"
}
```
Response:
```json
{
  "access": "<JWT_ACCESS_TOKEN>",
  "refresh": "<JWT_REFRESH_TOKEN>",
  "message": "Logged in successfully"
}
```

---

### Blogs

| Method | URL                        | Auth         | Description                       |
|--------|----------------------------|--------------|-----------------------------------|
| POST   | /blogs                     | Yes          | Create a new blog                 |
| GET    | /blogs                     | Yes          | List all blogs                    |
| GET    | /blogs/:id                 | Yes          | Get a blog by ID                  |
| DELETE | /blogs/:id                 | Owner/Admin  | Delete a blog                     |
| PATCH  | /blogs/:id                 | Owner        | Update a blog (partial)           |
| GET    | /blogs/paginated           | Yes          | Get paginated blogs               |
| GET    | /blogs/search              | Yes          | Search blogs by title/content     |
| GET    | /blogs/filter              | Yes          | Filter blogs by title/user with limit/offset |
| POST   | /blogs/:id/view            | Yes          | Increment blog view count         |
| POST   | /blogs/:id/like            | Yes          | Like a blog                       |
| DELETE | /blogs/:id/like            | Yes          | Unlike a blog                     |
| GET    | /blogs/:id/popularity      | Yes          | Get blog popularity (views, likes)|
| POST   | /blogs/ideas               | Yes          | Generate blog ideas (AI)          |
| POST   | /blogs/improve             | Yes          | Suggest blog improvements (AI)    |
| POST   | /blogs/:id/comments        | Yes          | Add a comment to a blog           |
| GET    | /blogs/:id/comments        | Yes          | List comments for a blog          |

#### Example: Create Blog
Request:
```json
{
  "title": "How to Use Go",
  "content": "Go is a statically typed, compiled language...",
  "tags": "go,programming,backend"
}
```
Response:
```json
{
  "message": "Blog created successfully",
  "blog": {
    "id": 1,
    "title": "How to Use Go",
    "content": "Go is a statically typed, compiled language...",
    "tags": ["go", "programming", "backend"]
  }
}
```

#### Example: Update Blog
Request (PATCH /blogs/:id):
```json
{
  "title": "Updated Title",
  "content": "Updated content"
}
```
Responses:
- 200: `{ "message": "blog updated" }`
- 404: `{ "error": "blog not found" }`

#### Example: Paginated Blogs
Request: GET /blogs/paginated?page=1&limit=10
Response:
```json
{
  "data": [ { "id": 1, "title": "..." } ],
  "total": 42,
  "page": 1,
  "limit": 10,
  "total_pages": 5
}
```

#### Example: Search Blogs
Request: GET /blogs/search?q=go&page=1&limit=10
Responses:
- 200:
```json
{
  "blogs": [ { "id": 1, "title": "Go Tips" } ],
  "meta": { "total": 5, "page": 1, "limit": 10 }
}
```
- 400 when q missing: `{ "error": "q is required" }`

#### Example: Track View / Like / Unlike / Popularity
- Track view: POST /blogs/1/view → `{ "message": "view tracked" }`
- Like: POST /blogs/1/like → `{ "message": "liked" }`
- Unlike: DELETE /blogs/1/like → `{ "message": "unliked" }`
- Popularity: GET /blogs/1/popularity → `{ "view_count": 12, "likes": 3 }`

---

### Comments

| Method | URL                   | Auth | Description                         |
|--------|------------------------|------|-------------------------------------|
| POST   | /blogs/:id/comments    | Yes  | Add a comment to a specific blog    |
| GET    | /blogs/:id/comments    | Yes  | List comments (paginated) for a blog|

Auth: Yes (Authorization header required)

Add Comment request body:
```json
{ "content": "Nice article!" }
```
Response 201:
```json
{
  "comment": {
    "id": 10,
    "content": "Nice article!",
    "user_id": 1,
    "blog_id": 2,
    "created_at": "2025-08-12T10:00:00Z"
  }
}
```

List Comments request: GET /blogs/2/comments?page=1&limit=10
Response 200:
```json
{
  "comments": [
    { "id": 10, "content": "Nice article!", "user_id": 1, "blog_id": 2 },
    { "id": 11, "content": "Great read", "user_id": 3, "blog_id": 2 }
  ],
  "meta": { "total": 2, "page": 1, "limit": 10 }
}
```

---

## Authentication and Authorization

- **Mechanism:** JWT tokens (access and refresh)
- **Login:** Returns access and refresh tokens
- **Protected Routes:** Require `Authorization: Bearer <token>` header
- **Roles:**
  - **User:** Can manage own profile, blogs, comments
  - **Admin:** Can promote/demote users, delete any blog
- **Middleware:** Enforces authentication, role, and ownership

---

## Data Models / Database Schema

### User
- id (int64, PK)
- username (string, unique)
- email (string, unique)
- password (string, hashed)
- role (string: user/admin)
- bio, profile_picture, phone, status
- created_at, updated_at

### Blog
- id (int64, PK)
- title, content
- user_id (FK to User)
- view_count, likes, dislikes
- created_at, updated_at

### Tag
- id (int64, PK)
- name (string, unique)
- content

### Comment
- id (int64, PK)
- content
- user_id (FK to User)
- blog_id (FK to Blog)
- created_at, updated_at

### Token
- id (int64, PK)
- type (access/refresh)
- content (JWT string)
- status (active/blocked)
- user_id (FK to User)

#### Relationships
- User 1--* Blog
- Blog *--* Tag (via join table)
- Blog 1--* Comment
- User 1--* Comment

---

## Business Logic / Use Cases

- **User Registration:** Validates input, hashes password, creates user, sends activation email
- **Login:** Validates credentials, issues JWT tokens
- **Profile Update:** Only owner or admin can update
- **Blog CRUD:** Authenticated users can create, update, delete their blogs; admins can delete any blog
- **Tag Management:** Tags are created/linked on blog creation
- **Password Reset:** Via email token or while logged in
- **Admin Actions:** Promote/demote users

---

## Middleware and Services

- **Auth Middleware:** Validates JWT, sets user context
- **Admin Middleware:** Checks for admin role
- **Account Owner Middleware:** Ensures user is acting on own resource
- **Error Handling Middleware:** Returns JSON errors
- **Email Service:** Sends activation and password reset emails (SMTP)
- **AI Service:** (Optional) Suggests blog ideas/improvements

---

## Error Handling

- All errors returned as JSON: `{ "error": "message" }`
- Uses appropriate HTTP status codes (400, 401, 403, 404, 500)

---

## Testing

- **Unit Tests:** For usecases, repositories, infrastructure
- **Integration Tests:** API happy-path flows (login, create blog, list, etc.)
- **Location:** `test/` directory
- **Run all tests:**
  ```sh
  go test ./...
  ```

Optional targeted runs:
```sh
go test ./test/usecases -v
go test ./test/repositories -v
```

---

## Setup and Deployment

### Prerequisites
- Go 1.20+
- PostgreSQL

### Environment Variables
Create a `.env` file:
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
OPENAI_API_KEY=your_openai_api_key_optional
```

### Running Locally
1. Clone the repo and install dependencies
2. Set up PostgreSQL and `.env`
3. Run migrations (auto-migrated on start)
4. Start server:
   ```sh
   go run delivery/main.go
   ```

Login and use the API:
1. Register: POST /register
2. Login: POST /login → copy `access` token
3. Call protected routes with header `Authorization: Bearer <access>`

### Deployment
- Build: `go build -o blog-platform delivery/main.go`
- Deploy binary and `.env` to server

---

## Security Considerations

- Passwords hashed with bcrypt
- JWT tokens with expiration and revocation
- Role-based access control
- Input validation and error handling
- Sensitive config via environment variables

---

## Future Improvements

- Comment moderation and editing
- Rate limiting and request logging
- API documentation (Swagger/OpenAPI)
- More granular permissions/roles
- WebSocket support for real-time updates

---

### Filter Blogs

Filter blogs using query parameters. All fields are optional unless stated otherwise.

- Endpoint: GET /blogs/filter
- Auth: Yes (Authorization: Bearer <token>)
- Query parameters:
  - title: string (optional, substring match on title)
  - user_id: int64 (optional, filter by author)
  - limit: int (optional, default 10, min 1)
  - offset: int (optional, default 0, min 0)

Example request (filter by title contains "go"):

GET /blogs/filter?title=go&limit=5&offset=0

Example 200 response:
[
  {
    "id": 42,
    "title": "Go Concurrency Patterns",
    "content": "…",
    "user_id": 123,
    "tags": [
      { "id": 1, "name": "go" },
      { "id": 2, "name": "concurrency" }
    ],
    "created_at": "2025-08-11T10:00:00Z",
    "updated_at": "2025-08-11T10:00:00Z"
  }
]

Example request (filter by user):

GET /blogs/filter?user_id=123&limit=10&offset=0

Example 400 response (invalid user_id):
{ "error": "invalid user_id" }

---