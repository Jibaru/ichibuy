-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    data JSONB,
    "timestamp" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_timestamp ON events("timestamp");
