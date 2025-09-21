-- +goose Up
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_stores_user_id ON stores(user_id);
CREATE INDEX idx_stores_name ON stores(name);
CREATE INDEX idx_stores_created_at ON stores(created_at);