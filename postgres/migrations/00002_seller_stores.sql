-- +goose Up
CREATE TABLE IF NOT EXISTS seller_stores (
    id              INT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,     
    public_id       TEXT         NOT NULL, -- create from name, natural key, behaves like a slug
    seller_id       BIGINT       NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    name            CITEXT       NOT NULL CHECK (LENGTH(name) <= 100),    -- TODO check if name is inserted uniquely   
    description     TEXT         NOT NULL,               
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ          
);

CREATE UNIQUE INDEX idx_seller_stores_name ON seller_stores(name);

CREATE INDEX idx_seller_stores_public_id ON seller_stores(public_id);

-- +goose Down
DROP TABLE IF EXISTS seller_stores;
