-- +goose Up
CREATE TYPE product_status AS ENUM (
    'Draft',
    'PendingApproval',
    'Published'
);

CREATE TABLE IF NOT EXISTS products (
    id                 BIGINT          NOT NULL DEFAULT generate_id() PRIMARY KEY,
    store_id           INT             NOT NULL REFERENCES stores(id) ON DELETE NO ACTION,
    category_id        INT             NOT NULL REFERENCES categories(id) ON DELETE NO ACTION,
    name               TEXT            NOT NULL CHECK (LENGTH(name) <= 255),
    brand              TEXT            CHECK (LENGTH(brand) <= 255),
    handle             TEXT            NOT NULL CHECK (LENGTH(handle) <= 120),
    image_urls         JSONB           NOT NULL,
    short_information  TEXT[],
    specifications     JSONB           NOT NULL,       -- BUG: storing specs in jsonb does not keep order.
    status             product_status  NOT NULL,
    created_at         TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ     NOT NULL DEFAULT NOW()

    -- TODO: possibly product_type column is needed
);

CREATE TABLE IF NOT EXISTS product_pricing (
    product_id          BIGINT          PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    base_price          NUMERIC         NOT NULL,
    currency_code       TEXT            NOT NULL CHECK (LENGTH(currency_code) <= 3),
    discount            NUMERIC         NOT NULL DEFAULT 0,
    discounted_price    NUMERIC         NOT NULL,
    is_on_sale          BOOLEAN         NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_variants (
    id                    BIGINT      NOT NULL DEFAULT generate_id() PRIMARY KEY,
    product_id            BIGINT      NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku                   TEXT        NOT NULL UNIQUE,  -- TODO: generate better sku values
    qty                   INTEGER     NOT NULL DEFAULT 0,
    is_custom_priced      BOOLEAN     NOT NULL DEFAULT false,
    image_identifiers     SMALLINT[],
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_variant_pricing (
    variant_id          BIGINT          PRIMARY KEY REFERENCES product_variants(id) ON DELETE CASCADE,
    base_price          NUMERIC         NOT NULL,
    currency_code       TEXT            NOT NULL CHECK (LENGTH(currency_code) <= 3),
    discount            NUMERIC         NOT NULL DEFAULT 0,
    discounted_price    NUMERIC         NOT NULL,
    is_on_sale          BOOLEAN         NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_attributes (
    id                  BIGINT  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id          BIGINT  NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute           TEXT    NOT NULL CHECK (LENGTH(attribute) <= 50),
    UNIQUE (product_id, attribute)
);

CREATE TABLE IF NOT EXISTS product_variant_attributes (
    variant_id              BIGINT   NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    product_attribute_id    BIGINT   NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    attribute_value         TEXT     NOT NULL CHECK (LENGTH(attribute_value) <= 50),
    
    PRIMARY KEY (variant_id, product_attribute_id)
);


-- +goose Down
DROP TABLE IF EXISTS product_variant_attributes;
DROP TABLE IF EXISTS product_attributes;
DROP TABLE IF EXISTS product_variant_pricing;
DROP TABLE IF EXISTS product_variants;
DROP TYPE IF EXISTS product_pricing;
DROP TYPE IF EXISTS products;
DROP TYPE IF EXISTS product_status;
