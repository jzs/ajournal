
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE Profile(
	UserID bigint references _User(id) NOT NULL UNIQUE,
	Name text NOT NULL,
	Email text NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE Profile;
