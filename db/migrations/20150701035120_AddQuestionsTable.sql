
-- +goose Up
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  full_text text NOT NULL,
  positions integer ARRAY
);


-- +goose Down
DROP TABLE questions;

