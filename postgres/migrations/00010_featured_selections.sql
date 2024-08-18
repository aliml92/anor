-- +goose Up
CREATE TABLE IF NOT EXISTS featured_selections (
   id              BIGINT NOT NULL DEFAULT generate_id() PRIMARY KEY,
   resource_path   TEXT NOT NULL CHECK (LENGTH(resource_path) <= 300),
   banner_info     JSON NOT NULL,
   image_url       TEXT NOT NULL,
   query_params    JSON,
   start_date      DATE,
   end_date        DATE,
   display_order   INT,
   created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
   updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS featured_selections;
