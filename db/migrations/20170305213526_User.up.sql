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
