-- +goose Up
CREATE TYPE cart_status AS ENUM (
    'Active',
    'Checkout',
    'Completed'
);

CREATE TABLE IF NOT EXISTS carts (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    user_id         BIGINT       REFERENCES users(id) ON DELETE CASCADE,
    status          cart_status  NOT NULL DEFAULT 'Active',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS cart_items (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    cart_id         BIGINT       NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    variant_id      BIGINT       NOT NULL REFERENCES product_variants(id) ON DELETE NO ACTION,
    qty             INTEGER      NOT NULL CHECK (qty > 0),
    price           NUMERIC      NOT NULL,
    thumbnail       TEXT         NOT NULL,
    product_name    TEXT         NOT NULL,
    variant_attributes  JSON,
    is_removed      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

