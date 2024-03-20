-- +goose Up
CREATE TYPE product_status AS ENUM (
    'Draft',
    'PendingApproval',
    'Active'    
);

CREATE TABLE IF NOT EXISTS products (
    id              BIGINT          NOT NULL DEFAULT generate_id() PRIMARY KEY, 
    store_id        INT             NOT NULL REFERENCES seller_stores(id) ON DELETE NO ACTION,
    category_id     INT             NOT NULL REFERENCES categories(id) ON DELETE NO ACTION,
    name            TEXT            NOT NULL CHECK (LENGTH(name) <= 255),
    brand           TEXT            CHECK (LENGTH(brand) <= 255),
    slug            TEXT            NOT NULL CHECK (LENGTH(slug) <= 120),
    short_info      TEXT[],
    image_urls      JSONB           NOT NULL,
    specs           JSONB           NOT NULL,
    status          product_status  NOT NULL,  
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ      
);

CREATE TABLE IF NOT EXISTS product_pricing (
    product_id          BIGINT          PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    base_price          NUMERIC         NOT NULL,
    currency_code       TEXT            NOT NULL CHECK (LENGTH(currency_code) <= 3),
    discount_level      NUMERIC         NOT NULL DEFAULT 0,
    discounted_amount   NUMERIC         NOT NULL DEFAULT 0,
    is_on_sale          BOOLEAN         NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS skus (
    id                  BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id          BIGINT      NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku                 TEXT        NOT NULL UNIQUE,
    quantity_in_stock   INTEGER     NOT NULL DEFAULT 0,
    has_sep_pricing     BOOLEAN     NOT NULL DEFAULT false,
    image_refs          SMALLINT[],  
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ    
);

CREATE TABLE IF NOT EXISTS sku_pricing (
    sku_id              BIGINT          PRIMARY KEY REFERENCES skus(id) ON DELETE CASCADE,
    base_price          NUMERIC         NOT NULL,
    currency_code       TEXT            NOT NULL CHECK (LENGTH(currency_code) <= 3),
    discount_level      NUMERIC         NOT NULL DEFAULT 0,
    discounted_amount   NUMERIC         NOT NULL DEFAULT 0,
    is_on_sale          BOOLEAN         NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_attributes (
    id                  BIGINT  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id          BIGINT  NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute           TEXT    NOT NULL CHECK (LENGTH(attribute) <= 50),
    UNIQUE (product_id, attribute)
);

CREATE TABLE IF NOT EXISTS sku_product_attribute_values (
    sku_id                  BIGINT   NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
    product_attribute_id    BIGINT   NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    attribute_value         TEXT     NOT NULL CHECK (LENGTH(attribute_value) <= 50),
    
    PRIMARY KEY (sku_id, product_attribute_id)
);