-- +goose Up
CREATE TABLE IF NOT EXISTS coupons (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(50)  NOT NULL UNIQUE, -- Unique code for the coupon/promo
    description     TEXT,                         -- Description of the coupon/promo
    discount_type   VARCHAR(20)  NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')), -- Type of discount (percentage or fixed amount)
    discount_value  DECIMAL(10, 2) NOT NULL,      -- Value of the discount
    max_uses        INT,                          -- Maximum number of times the coupon can be used (optional)
    expiration_date DATE,                         -- Expiration date of the coupon (optional)
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupon_usage (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
    coupon_code     VARCHAR(50)  NOT NULL, -- Applied coupon code
    user_id         BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- User who applied the coupon
    cart_id         BIGINT       NOT NULL REFERENCES carts(id) ON DELETE CASCADE, -- Cart where the coupon was applied
    applied_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW() -- Timestamp when the coupon was applied
);
