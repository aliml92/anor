-- +goose Up
CREATE TYPE order_status AS ENUM (
    'Pending',
    'Processing',
    'Shipped',
    'Delivered',
    'Canceled'
);

CREATE TABLE IF NOT EXISTS orders (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    cart_id         BIGINT       NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    user_id         BIGINT       REFERENCES users(id) ON DELETE NO ACTION,
    status          order_status NOT NULL DEFAULT 'Pending',
    total_amount    NUMERIC      NOT NULL,
    payment_intent_id TEXT       NOT NULL,
    shipping_address JSON        NOT NULL,
    billing_address  JSON        NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    order_id        BIGINT       NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id      BIGINT       NOT NULL REFERENCES product_variants(id) ON DELETE NO ACTION,
    qty             INTEGER      NOT NULL CHECK (qty > 0),
    price           NUMERIC      NOT NULL,
    thumbnail       TEXT         NOT NULL,
    product_name    TEXT         NOT NULL,
    variant_attributes JSON,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
