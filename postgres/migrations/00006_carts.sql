-- +goose Up
DROP TYPE IF EXISTS cart_status CASCADE;
CREATE TYPE cart_status AS ENUM (
    'Open',
    'Completed',
    'Merged',
    'Expired',
    'Abandoned'
);

CREATE TABLE IF NOT EXISTS carts (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    user_id         BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    status          cart_status  NOT NULL DEFAULT 'Open',
    last_activity   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);

CREATE TABLE IF NOT EXISTS cart_items (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    cart_id         BIGINT       NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    variant_id      BIGINT       NOT NULL REFERENCES product_variants(id) ON DELETE NO ACTION,
    qty             INTEGER      NOT NULL CHECK (qty > 0),
    price           NUMERIC      NOT NULL,
    currency        TEXT         NOT NULL DEFAULT 'USD' CHECK (LENGTH(currency) <= 3),
    thumbnail       TEXT         NOT NULL,
    product_name    TEXT         NOT NULL,
    product_path    TEXT         NOT NULL,
    variant_attributes  JSON,
    is_removed      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_variant_id ON cart_items(variant_id);

ALTER TABLE cart_items ADD CONSTRAINT  unique_cart_variant UNIQUE (cart_id, variant_id);


-- +goose Down
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TYPE IF EXISTS cart_status;