
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

CREATE TABLE Entry(
	ID BIGSERIAL primary key NOT NULL,
	JournalID bigint references Journal(id) NOT NULL,
	Date timestamp NOT NULL,
	Title text NOT NULL,
	Content text NOT NULL,
	Created timestamp NOT NULL,
	Published timestamp,
	IsPublished boolean NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE Entry;
DROP TABLE Journal;
