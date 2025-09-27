-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    store_id UUID NOT NULL,
    images JSONB,
    prices JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_store FOREIGN KEY(store_id) REFERENCES stores(id) ON DELETE CASCADE
);

CREATE INDEX idx_products_store_id ON products(store_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_created_at ON products(created_at);
