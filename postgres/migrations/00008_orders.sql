-- +goose Up
DROP TYPE IF EXISTS payment_status CASCADE;
CREATE TYPE payment_status AS ENUM (
    'Pending',   -- Payment attempted but failed or payment is being processed or will be paid on delivery.
    'Paid'       -- Paid successfully
);

DROP TYPE IF EXISTS order_status CASCADE;
CREATE TYPE order_status AS ENUM (
    'Pending',    -- Order received, not yet processed
    'Processing', -- Order is being prepared
    'Shipped',    -- Order has been shipped
    'Delivered',  -- Order has been delivered
    'Cancelled',  -- Order was cancelled
    'Fulfilled'   -- Customer received the order
);

DROP TYPE IF EXISTS payment_method CASCADE;
CREATE TYPE payment_method AS ENUM ('StripeCard', 'Paypal', 'AnorInstallment', 'PayOnDelivery');

CREATE TABLE IF NOT EXISTS orders (
  id                  BIGINT          NOT NULL DEFAULT generate_id() PRIMARY KEY,
  user_id             BIGINT          NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
  cart_id             BIGINT          NOT NULL REFERENCES carts(id),
  payment_method      payment_method  NOT NULL,
  payment_status      payment_status  NOT NULL DEFAULT 'Pending',
  status              order_status    NOT NULL DEFAULT 'Pending',
  shipping_address_id BIGINT          NOT NULL REFERENCES addresses(id),
  is_pickup           BOOLEAN         NOT NULL DEFAULT false,
  amount              NUMERIC         NOT NULL,
  currency            CHAR(3)         NOT NULL,
  created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_cart_id ON orders(cart_id);

CREATE TABLE IF NOT EXISTS order_items (
   id                  BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY,
   order_id            BIGINT       NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
   variant_id          BIGINT       NOT NULL REFERENCES product_variants(id) ON DELETE NO ACTION,
   qty                 INTEGER      NOT NULL CHECK (qty > 0),
   price               NUMERIC      NOT NULL,
   thumbnail           TEXT         NOT NULL,
   product_name        TEXT         NOT NULL,
   variant_attributes  JSONB,
   created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
   updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

-- +goose Down
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS payment_method;