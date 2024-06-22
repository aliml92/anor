-- +goose Up
CREATE TABLE IF NOT EXISTS wishlists (
    id           BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    user_id      BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS wishlist_items (
    id            BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    wishlist_id   BIGINT       NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    variant_id    BIGINT       NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    price         NUMERIC      NOT NULL,
    thumbnail     TEXT         NOT NULL,
    product_name  TEXT         NOT NULL,
    variant_attributes  JSON,
    is_removed    BOOLEAN      NOT NULL DEFAULT FALSE,
    added_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (wishlist_id, variant_id)
);

CREATE TABLE IF NOT EXISTS wishlist_shares (
    id                    BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    wishlist_id           BIGINT       NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    share_token           TEXT         NOT NULL UNIQUE,
    share_token_expiration TIMESTAMPTZ NOT NULL,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
