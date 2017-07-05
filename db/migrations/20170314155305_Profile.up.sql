CREATE TABLE Profile(
	UserID bigint references _User(id) NOT NULL UNIQUE,
	Name text NOT NULL,
	Email text NOT NULL
);

CREATE TABLE Subscription(
	UserID bigint references _User(id) NOT NULL UNIQUE,
	StripeCustomerID text NOT NULL,
	StripeSubscriptionID text
);
