-- +goose Up
CREATE TABLE stripe_card_payments (
     id                       BIGINT  NOT NULL DEFAULT generate_id() PRIMARY KEY,
     order_id                 BIGINT  NOT NULL REFERENCES orders(id),
     user_id                  BIGINT  REFERENCES users(id),
     billing_address_id       BIGINT  NOT NULL REFERENCES addresses(id),
     payment_intent_id        TEXT    NOT NULL,
     payment_method_id        TEXT,
     amount                   NUMERIC NOT NULL,
     currency                 CHAR(3) NOT NULL,
     status                   TEXT    NOT NULL,
     client_secret            TEXT,
     last_error               TEXT,
     card_last4               CHAR(4) NOT NULL,
     card_brand               TEXT    NOT NULL,
     created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PayPal payments
CREATE TABLE paypal_payments (
     id                   BIGINT   NOT NULL DEFAULT generate_id() PRIMARY KEY,
     order_id             BIGINT   NOT NULL REFERENCES orders(id),
     user_id              BIGINT   REFERENCES users(id),
     billing_address_id   BIGINT   NOT NULL REFERENCES addresses(id),
     paypal_order_id      TEXT     NOT NULL,
     paypal_payer_id      TEXT,
     amount               NUMERIC  NOT NULL,
     currency             CHAR(3)  NOT NULL,
     status               TEXT     NOT NULL,
     last_error           TEXT,
     created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Installment payments
CREATE TABLE anor_installment_payments (
    id                  BIGINT  NOT NULL DEFAULT generate_id() PRIMARY KEY,
    order_id            BIGINT  NOT NULL REFERENCES orders(id),
    user_id             BIGINT  REFERENCES users(id),
    billing_address_id  BIGINT  NOT NULL REFERENCES addresses(id),
    provider            TEXT    NOT NULL,
    anor_installment_plan_id TEXT    NOT NULL,
    total_amount        NUMERIC NOT NULL,
    currency            CHAR(3) NOT NULL,
    number_of_installments INTEGER NOT NULL,
    status              TEXT NOT NULL,
    last_error          TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Pay on Delivery/Pickup
CREATE TABLE pay_on_delivery_payments (
    id                 BIGINT   NOT NULL DEFAULT generate_id() PRIMARY KEY,
    order_id           BIGINT   NOT NULL REFERENCES orders(id),
    user_id            BIGINT   REFERENCES users(id),
    billing_address_id BIGINT   NOT NULL REFERENCES addresses(id),
    amount             NUMERIC  NOT NULL,
    currency           CHAR(3)  NOT NULL,
    payment_type       TEXT     NOT NULL CHECK (payment_type IN ('Cash', 'QRCode', 'Card')),
    status             TEXT     NOT NULL,
    tracking_number    TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS pay_on_delivery_payments;
DROP TABLE IF EXISTS anor_installment_payments;
DROP TABLE IF EXISTS paypal_payments;
DROP TABLE IF EXISTS stripe_card_payments;
