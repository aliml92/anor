-- +goose Up
ALTER TABLE carts RENAME COLUMN stripe_client_secret TO pi_client_secret;

