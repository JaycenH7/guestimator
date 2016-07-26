
-- +goose Up
ALTER TABLE wikipages RENAME COLUMN page_id TO id;


-- +goose Down
ALTER TABLE wikipages RENAME COLUMN id TO page_id;

