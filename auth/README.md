# ichibuy/auth

OAuth authentication microservice for the ichibuy platform.

## Overview

This service provides OAuth authentication (Google) with JWT token generation using RSA-256 signing. It exposes a JWKS endpoint for public key distribution and JWT verification.

## Features

- Google OAuth 2.0 authentication
- RSA-256 JWT token generation and signing
- JWKS endpoint for public key distribution
- PostgreSQL user persistence
- Swagger API documentation
- Vercel serverless deployment support

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/auth/google` | Initiate Google OAuth flow |
| GET | `/api/v1/auth/google/callback` | Handle OAuth callback |
| GET | `/api/v1/auth/.well-known/jwks.json` | JSON Web Key Set |
| GET | `/api/swagger/*` | API documentation |

## Configuration

Required environment variables:

```bash
POSTGRES_URI=postgresql://user:pass@host:port/database
API_PORT=8080
API_BASE_URI=https://your-api-domain.com
WEB_BASE_URI=https://your-web-domain.com
JWT_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----..."
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

## Setup

1. **Generate RSA keys:**
   ```bash
   go run scripts/generate_rsa_keys.go
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

3. **Run database migrations:**
   ```bash
   make migrate-up
   ```

4. **Start the server:**
   ```bash
   make run
   ```

## Development

- **Build:** `make build`
- **Generate docs:** `make build` (includes swagger generation)
- **Generate DAOs:** `make gen`

## Database

Uses PostgreSQL with Goose migrations. Migration files are located in `./db/migrations/`.

## Deployment

Configured for Vercel serverless deployment. The `vercel.json` file contains the necessary configuration.

## Authentication Flow

1. User initiates OAuth via `/auth/google`
2. User completes OAuth on Google's platform
3. Google redirects to `/auth/google/callback`
4. Service creates/finds user and generates JWT
5. User is redirected to web application with JWT

## JWT Structure

Generated JWTs include:
- `user_id`: User identifier
- `email`: User email address
- `exp`: Token expiration
- `iat`: Issued at timestamp
- `iss`: Issuer (ichibuy-auth)
- `kid`: Key identifier for verification

## Security

- Uses RSA-256 for JWT signing
- Private key remains secure on server
- Public key exposed via JWKS endpoint
- CORS configured for cross-origin requests
