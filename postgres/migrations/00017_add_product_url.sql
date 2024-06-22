-- +goose Up
ALTER TABLE cart_items ADD COLUMN product_path TEXT NOT NULL ;
