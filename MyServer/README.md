# magicgate

`magicgate` is a Go-based REST API server designed to manage users, cryptographic keys, and perform encryption/decryption operations. It uses PostgreSQL as its database and JWT for authentication.

## Features

- **User Management**: Create, retrieve, update, and delete users.
- **Authentication**: User login with username/password, generating a JSON Web Token (JWT).
- **Key Management**: Create, retrieve, update, and delete cryptographic keys associated with users. Keys are stored securely (as `BYTEA` in DB, not exposed via API).
- **Encryption/Decryption**: API endpoints to encrypt and decrypt data using a user's stored keys and Go's `crypto` package (AES-256 GCM).
- **PostgreSQL Database**: Persistent storage for users and keys.
- **Secure Passwords**: User passwords are hashed using bcrypt.
- **JWT Authentication Middleware**: Protects key management and crypto endpoints.

## Project Structure

```
magicgate/MyServer
├── main.go               # Main application entry point
├── go.mod                # Go module file
├── go.sum                # Go module checksums
├── config/
│   └── config.go         # Application configuration loading (from .env or env vars)
├── database/
│   ├── db.go             # Database connection and table creation
│   ├── models.go         # Database models (User, Key)
│   ├── user_repo.go      # CRUD operations for User
│   └── key_repo.go       # CRUD operations for Key
├── handlers/
│   ├── user_handlers.go  # HTTP handlers for User CRUD
│   ├── key_handlers.go   # HTTP handlers for Key CRUD
│   ├── auth_handlers.go  # HTTP handler for Login (JWT generation)
│   └── crypto_handlers.go# HTTP handlers for Encryption/Decryption
├── middleware/
│   └── auth_middleware.go# JWT authentication middleware
└── utils/
    ├── jwt.go            # JWT token generation and validation
    ├── password.go       # Password hashing and comparison
    └── crypto.go         # Cryptographic utility functions (AES-GCM)
```

## Getting Started

### Prerequisites

- Go (version 1.21 or higher)
- PostgreSQL database
- `make` (optional, for convenience)

### 1. Setup PostgreSQL

Ensure you have a PostgreSQL instance running. Create a database, e.g., `magicgate`.

```sql
CREATE DATABASE magicgate;
CREATE USER user WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE magicgate TO user;
```

### 2. Configuration

Create a `.env` file in the root directory of the project:

```
DATABASE_URL="postgres://user:password@localhost:5432/magicgate?sslmode=disable"
JWT_SECRET="your_super_secret_jwt_key_here" # **IMPORTANT: Change this to a strong, unique key!**
SERVER_PORT="8080"
ENCRYPTION_NONCE_SIZE="12" # Recommended GCM nonce size
```

Replace `user`, `password`, `localhost:5432`, and `magicgate` with your PostgreSQL credentials and connection details.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

The server will start on the port specified in `SERVER_PORT` (default: `8080`).

## API Endpoints

### User Endpoints (Public)

- `POST /register`: Create a new user.
- `POST /login`: Authenticate a user and get a JWT token.

### Authenticated Endpoints (Require `Authorization: Bearer <JWT_TOKEN>` header)

- **User CRUD** (for simplicity, any authenticated user can access any user ID for now; a more secure implementation would restrict access to the user's own profile or admin roles):
    - `GET /api/users`: Get all users.
    - `GET /api/users/{id}`: Get a user by ID.
    - `PUT /api/users/{id}`: Update a user by ID.
    - `DELETE /api/users/{id}`: Delete a user by ID.
- **Key CRUD** (user-specific):
    - `POST /api/keys`: Create a new cryptographic key for the authenticated user.
    - `GET /api/keys`: Get all keys for the authenticated user.
    - `GET /api/keys/{id}`: Get a specific key for the authenticated user.
    - `PUT /api/keys/{id}`: Update a key's name for the authenticated user.
    - `DELETE /api/keys/{id}`: Delete a key for the authenticated user.
- **Crypto Operations** (user-specific):
    - `POST /api/encrypt`: Encrypt data using a specified key owned by the authenticated user.
    - `POST /api/decrypt`: Decrypt data using a specified key owned by the authenticated user.

## Example Usage (using `curl`)

### 1. Register a user

```bash
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username": "testuser", "password": "password123"}'
```

### 2. Login and get JWT

```bash
TOKEN=$(curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username": "testuser", "password": "password123"}' | jq -r .token)
echo "JWT Token: $TOKEN"
```

### 3. Create a key

```bash
KEY_ID=$(curl -X POST http://localhost:8080/api/keys -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name": "my_first_key"}' | jq -r .id)
echo "Created Key ID: $KEY_ID"
```

### 4. Get all keys

```bash
curl -X GET http://localhost:8080/api/keys -H "Authorization: Bearer $TOKEN"
```

### 5. Encrypt data

```bash
ENCRYPTED_PAYLOAD=$(curl -X POST http://localhost:8080/api/encrypt -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "{\"key_id\": $KEY_ID, \"payload\": \"This is my secret message.\"}" | jq -r .result)
echo "Encrypted Payload: $ENCRYPTED_PAYLOAD"
```

### 6. Decrypt data

```bash
curl -X POST http://localhost:8080/api/decrypt -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "{\"key_id\": $KEY_ID, \"payload\": \"$ENCRYPTED_PAYLOAD\"}"
```

### 7. Delete a key

```bash
curl -X DELETE http://localhost:8080/api/keys/$KEY_ID -H "Authorization: Bearer $TOKEN" -v
```

## Security Considerations

- **JWT Secret**: The `JWT_SECRET` in `.env` should be a strong, randomly generated string and kept confidential.
- **Key Management**: Storing raw cryptographic key material directly in the database, even as `BYTEA`, is generally not recommended for high-security applications. A more robust solution would involve a Key Management System (KMS) or hardware security modules (HSMs). This example demonstrates the cryptographic operations but simplifies key storage.
- **Password Hashing**: Bcrypt is used, which is good.
- **Error Handling**: The error handling is basic. In a production system, more detailed logging and user-friendly error messages (without exposing internal details) would be needed.
- **Input Validation**: Input validation is minimal. Robust validation should be added for all API inputs.
- **HTTPS**: Always use HTTPS in production to protect data in transit.
- **User Authorization**: The current user CRUD endpoints are protected by JWT, but any authenticated user can technically access any user ID. For a real application, implement fine-grained authorization (e.g., a user can only access their own profile, or only admins can access all users).
