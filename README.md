# mwx-go-auth-service

[![CD](https://github.com/roymwxuk/mwx-go-auth-service/actions/workflows/cd.yml/badge.svg)](https://github.com/roymwxuk/mwx-go-auth-service/actions/workflows/cd.yml)

A Go authentication microservice supporting Google Sign-In and JWT authentication.

Features:
- Google Sign-In
- JWT authentication (RS384)
- Refresh tokens
- User identity mapping
- Designed for microservice architecture
- GitHub Actions CI

---

## Authentication Flow

```text
Google ID Token
       │
       ▼
Verify with Google
       │
       ▼
Find user by
(provider, provider_user_id)
       │
 ┌─────┴─────┐
 │           │
Found     Not Found
 │           │
 │       Create User
 │       Create Identity
 │           │
 └─────► Generate JWT
             │
             ▼
Return Access Token + Refresh Token
```

## Database schema

```text
users
-----
id (UUID)
email
display_name
avatar_url
role
status
created_at
updated_at

user_identities
---------------
id
user_id
provider
provider_user_id
created_at
updated_at
```

## APIs

GET /health

POST /auth/google
POST /auth/refresh

GET  /users/me

## Tech Stack

- Go
- Gin
- PostgreSQL
- pgx
- Goose
- Google Identity
- JWT (RS384)
