-- +goose Up
CREATE TABLE IF NOT EXISTS stores (
    id              INT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,     
    handle          TEXT         NOT NULL, -- create from name, natural key, behaves like a slug
    user_id         BIGINT       NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    name            CITEXT       NOT NULL CHECK (LENGTH(name) <= 100),    -- TODO check if name is inserted uniquely   
    description     TEXT         NOT NULL,               
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()

    -- TODO: additional columns for logo, banner image, address, phone number, legal information and so on.
);

-- TODO: add store_reviews, store_ratings, store_statistics if needed

CREATE UNIQUE INDEX idx_stores_name ON stores(name);

CREATE INDEX idx_stores_handle ON stores(handle);

-- +goose Down
DROP TABLE IF EXISTS stores;
