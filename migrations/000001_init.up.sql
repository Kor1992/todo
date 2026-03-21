CREATE SCHEMA IF NOT EXISTS todoapp;

CREATE TABLE todoapp.users (
    id           SERIAL       PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)           CHECK (phone_number ~ '^\+?[0-9]+$' AND (char_length(phone_number) BETWEEN 10 AND 15)
);