
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE Profile ADD COLUMN Description text NOT NULL DEFAULT('');


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE Profile DROP COLUMN Description;
