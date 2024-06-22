-- +goose Up

CREATE SEQUENCE global_id_sequence;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_id(OUT result BIGINT) AS $$
  DECLARE

    our_epoch BIGINT := 1483228800;

    seq_id BIGINT;
    now_millis BIGINT;
    service_id INT := 1;
  BEGIN
      SELECT nextval('global_id_sequence') % 1024 INTO seq_id;
      SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp())) INTO now_millis;
      result := (now_millis - our_epoch) << 20;
      result := result | (service_id << 10);
      result := result | (seq_id);
  END;
$$ LANGUAGE PLPGSQL;
-- +goose StatementEnd

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_status AS ENUM (
    'Blocked', 
    'RegistrationPending', 
    'Active'
);

CREATE TABLE IF NOT EXISTS users (
    id              BIGINT       NOT NULL DEFAULT generate_id() PRIMARY KEY, 
    email           CITEXT       NOT NULL UNIQUE,
    password        TEXT         NOT NULL CHECK (LENGTH(password) <= 255), -- User's password (hashed)
    phone_number    TEXT         CHECK (LENGTH(phone_number) <= 20),                
    full_name       TEXT         NOT NULL CHECK (LENGTH(full_name) <= 100),
    status          user_status  NOT NULL DEFAULT 'RegistrationPending',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS roles (
    id              SMALLINT     GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    role            TEXT         UNIQUE NOT NULL
);

ALTER TABLE roles
ADD CONSTRAINT chk_roles_role_length CHECK (LENGTH(role) <= 20);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id         BIGINT   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id         SMALLINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

INSERT INTO roles (role) VALUES 
    ( 'admin' ),
    ( 'seller' ),
    ( 'customer');

-- TODO: add permissions table if needed

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TYPE IF EXISTS user_status CASCADE;