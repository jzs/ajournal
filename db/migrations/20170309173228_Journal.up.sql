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
