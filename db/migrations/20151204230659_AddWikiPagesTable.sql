
-- +goose Up
CREATE TABLE wiki_pages (
  page_id integer PRIMARY KEY,
  title text NOT NULL
);

ALTER TABLE questions
ADD COLUMN page_id integer REFERENCES wiki_pages NOT NULL;


-- +goose Down
ALTER TABLE questions
DROP COLUMN page_id;

DROP TABLE wiki_pages;
