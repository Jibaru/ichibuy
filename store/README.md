# Store Microservice

A microservice for managing stores and customers in the ichibuy application.

## Features

- **Store Management**: CRUD operations for stores with name, description, location, and slug generation
- **Customer Management**: CRUD operations for customers with email and phone validation
- **JWT Authentication**: Validates JWT tokens from the auth microservice
- **GraphQL API**: Query stores with filtering, sorting, and pagination
- **Event Bus**: Publishes events for store and customer operations
- **Value Objects**: Email and phone validation using domain-driven design

## API Endpoints

### Stores
- `POST /api/v1/stores` - Create a new store
- `GET /api/v1/stores/:id` - Get store by ID
- `PUT /api/v1/stores/:id` - Update store
- `DELETE /api/v1/stores/:id` - Delete store
- `GET /api/v1/stores` - List stores with filters and pagination

### Customers
- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers/:id` - Get customer by ID
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

### Products
- `POST /api/v1/products` - Create a new product
- `GET /api/v1/products/:id` - Get product by ID
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `GET /api/v1/products` - List products with filters and pagination

### GraphQL
- `POST /api/v1/graphql` - GraphQL endpoint for querying stores

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