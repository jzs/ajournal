CREATE TABLE OauthUser (
	UserID bigint references _User(id) NOT NULL,
	Provider text NOT NULL,
	ProviderUsername text NOT NULL
);
