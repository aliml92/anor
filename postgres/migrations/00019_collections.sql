-- +goose Up
CREATE TABLE IF NOT EXISTS collections (
   id              INT  NOT NULL DEFAULT generate_category_id() PRIMARY KEY,
   name            TEXT NOT NULL CHECK (LENGTH(name) <= 100),
   handle          TEXT NOT NULL,
   description     TEXT
);

CREATE TABLE IF NOT EXISTS product_collections (
   product_id      BIGINT  NOT NULL,
   collection_id   INT     NOT NULL,
   PRIMARY KEY (product_id, collection_id),
   FOREIGN KEY (product_id)    REFERENCES products(id),
   FOREIGN KEY (collection_id) REFERENCES collections(id)
);