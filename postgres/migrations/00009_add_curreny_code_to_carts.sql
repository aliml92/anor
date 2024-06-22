-- +goose Up
ALTER TABLE cart_items
    ADD COLUMN currency_code TEXT NOT NULL DEFAULT 'USD' CHECK (LENGTH(currency_code) <= 3);

