
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE Journal(
	ID BIGSERIAL primary key NOT NULL,
	Title text NOT NULL,
	Description text NOT NULL,
	UserID bigint references _User(id) NOT NULL,
	Public boolean NOT NULL DEFAULT(false),
	Created timestamp NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE Journal;
