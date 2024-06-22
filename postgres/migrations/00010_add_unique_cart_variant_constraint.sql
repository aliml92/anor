-- +goose Up
ALTER TABLE cart_items ADD CONSTRAINT unique_cart_variant UNIQUE (cart_id, variant_id);
