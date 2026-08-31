# mwx-go-auth-service

[![CD](https://github.com/roymwxuk/mwx-go-auth-service/actions/workflows/cd.yml/badge.svg)](https://github.com/roymwxuk/mwx-go-auth-service/actions/workflows/cd.yml)

A Go authentication microservice supporting Google Sign-In, JWT authentication, and secure HttpOnly cookie-based sessions.

## Features

- Google Sign-In
- JWT authentication (RS384)
- HttpOnly cookie-based authentication
- Refresh token rotation
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
Set HttpOnly Cookies
(access_token, refresh_token)
       │
       ▼
Return User Profile
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

## Authentication

### Browser Clients

Authentication uses HttpOnly cookies.

After a successful login, the service issues:

- access_token
- refresh_token

Both tokens are stored as secure HttpOnly cookies and are automatically sent by the browser with subsequent requests.

---

## Roadmap

- Token endpoint for non-browser clients (`POST /auth/token`)
- MCP client authentication support
- Apple Sign-In
- Email & Password authentication
- OAuth 2.0 / OpenID Connect compatibility

---

## Tech Stack

- Go
- Gin
- PostgreSQL
- pgx
- Goose
- Google Identity
- JWT (RS384)
