-- +goose Up
CREATE TABLE IF NOT EXISTS featured_promotions (
   id              BIGINT NOT NULL DEFAULT generate_id() PRIMARY KEY,
   title           TEXT   NOT NULL CHECK (LENGTH(title) <= 100),
   image_url       TEXT   NOT NULL,
   type            TEXT   NOT NULL CHECK (type IN ('category', 'collection')),
   target_id       INT,
   filter_params   JSONB,
   start_date      DATE,
   end_date        DATE,
   display_order   INT
);