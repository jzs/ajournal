
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE _User(
	ID BIGSERIAL primary key NOT NULL,
	Username text NOT NULL UNIQUE,
	Password text NOT NULL,
	Active boolean NOT NULL DEFAULT(false),
	Created timestamp NOT NULL
);

CREATE TABLE UserToken(
	Token text NOT NULL UNIQUE,
	UserID bigint references _User(ID) NOT NULL,
	Expires timestamp NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE UserToken;
DROP TABLE _User;
