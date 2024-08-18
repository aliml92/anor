-- +goose Up
DROP TYPE IF EXISTS address_default_type CASCADE;
CREATE TYPE address_default_type AS ENUM (
    'Shipping',
    'Billing'
);

CREATE TABLE IF NOT EXISTS addresses (
     id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
     user_id         BIGINT       NOT NULL REFERENCES users(id),
     default_for     address_default_type,
     name            TEXT         NOT NULL CHECK (LENGTH(name) <= 100),
     address_line_1  TEXT         NOT NULL CHECK (LENGTH(address_line_1) <= 100),
     address_line_2  TEXT         CHECK (LENGTH(address_line_2) <= 100),
     city            TEXT         NOT NULL CHECK (LENGTH(city) <= 100),
     state_province  TEXT         CHECK (LENGTH(state_province) <= 100),
     postal_code     TEXT         CHECK (LENGTH(postal_code) <= 20),
     country         TEXT         CHECK (LENGTH(country) <= 100),
     phone           TEXT         CHECK (LENGTH(phone) <= 20),
     created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
     updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Create an index on user_id and address_type for faster queries
CREATE INDEX IF NOT EXISTS idx_addresses_user_id_type ON addresses (user_id, default_for);

-- Add a constraint to ensure a user has only one default address per type
ALTER TABLE addresses
    ADD CONSTRAINT unique_default_address
        UNIQUE (user_id, default_for);


-- +goose Down
DROP TABLE IF EXISTS addresses;
DROP TYPE IF EXISTS  address_default_type;