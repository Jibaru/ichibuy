# Order Microservice

A microservice for managing stores and customers in the ichibuy application.

## Features

- <complete>

## API Endpoints

<Complete>


## Environment Variables

Copy `.env.example` to `.env` and configure it.

## Running the Service

1. Install dependencies:
```bash
make deps
```

2. Run the service:
```bash
make run
```

## Database Setup

Run the migrations in the `db/migrations/` directory to set up the database schema.

## Authentication

All endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

The token must contain a `user_id` claim, which is validated against the auth microservice's JWKS endpoint.
