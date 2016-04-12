
-- +goose Up
ALTER TABLE wiki_pages RENAME TO wikipages;
ALTER TABLE questions RENAME COLUMN page_id TO wikipage_id;


-- +goose Down
ALTER TABLE questions RENAME COLUMN wiki_page_id TO page_id;
ALTER TABLE wikipages RENAME TO wiki_pages;

