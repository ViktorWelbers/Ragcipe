-- +goose Up
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    servings INTEGER NOT NULL,
    servings_type VARCHAR(50) NOT NULL,
    country_code VARCHAR(50) NOT NULL,
    host_url VARCHAR(50) NOT NULL,
    original_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
