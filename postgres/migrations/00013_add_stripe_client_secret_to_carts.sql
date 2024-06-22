-- +goose Up
ALTER TABLE carts ADD COLUMN stripe_client_secret TEXT;