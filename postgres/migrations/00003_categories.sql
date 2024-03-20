-- +goose Up
CREATE SEQUENCE category_id_sequence START 1000000 INCREMENT BY 10000;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_category_id(OUT result INT)
AS $$
DECLARE
	seq_id INT;
    random_offset INT;
	low SMALLINT := 100;
	high SMALLINT := 10000;
BEGIN
    -- Get the next value from the sequence
    seq_id := NEXTVAL('category_id_sequence');

    -- Generate a random offset within the specified range
    random_offset := floor(random()* (high-low + 1) + low)::INT;

    -- Add the random offset to the category ID
    result := seq_id + random_offset;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS categories (
    id              INT  NOT NULL DEFAULT generate_category_id() PRIMARY KEY,
    category        TEXT NOT NULL CHECK (LENGTH(category) <= 100),
    slug            TEXT NOT NULL, 
    parent_id       INT,
    
    UNIQUE(category, parent_id)
);


ALTER TABLE categories
ADD CONSTRAINT fk_parent_category
FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE NO ACTION;


-- +goose Down
ALTER TABLE IF EXISTS categories
DROP CONSTRAINT IF EXISTS fk_parent_category;

DROP TABLE IF EXISTS categories;