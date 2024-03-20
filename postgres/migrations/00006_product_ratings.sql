-- +goose Up
CREATE TABLE IF NOT EXISTS product_ratings (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY, 
    product_id      BIGINT       NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id         BIGINT       REFERENCES users(id) ON DELETE SET NULL, 
    rating          SMALLINT     NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review          TEXT         CHECK (LENGTH(review) <= 1000),
    image_urls      TEXT[],
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()        
);  

