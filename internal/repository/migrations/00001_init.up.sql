CREATE TABLE IF NOT EXISTS users (
	id serial NOT NULL UNIQUE,
	login varchar(64) NOT NULL UNIQUE,
	password varchar(64) NOT NULL,
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS users_id ON users (id);
CREATE INDEX IF NOT EXISTS users_login ON users (login);

CREATE TYPE order_status AS ENUM
	('NEW', 'REGISTERED', 'PROCESSING', 'PROCESSED' ,'INVALID');

CREATE TABLE IF NOT EXISTS orders (
	id bigint NOT NULL UNIQUE,
	user_id serial NOT NULL,
	status order_status NOT NULL DEFAULT 'NEW'::order_status,
	accrual money NOT NULL DEFAULT '0',
	sum money NOT NULL DEFAULT '0',
	date timestamp with time zone NOT NULL,
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS orders_id ON orders (id);
CREATE INDEX IF NOT EXISTS orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS orders_date ON orders (date);

CREATE TABLE IF NOT EXISTS accounts (
	id serial NOT NULL UNIQUE,
	user_id bigint NOT NULL UNIQUE,
	amount money NOT NULL DEFAULT '0',
	PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS accounts_user_id ON accounts (user_id);

ALTER TABLE Orders ADD CONSTRAINT Orders_fk0 FOREIGN KEY (user_id) REFERENCES Users(id);
ALTER TABLE Accounts ADD CONSTRAINT Accounts_fk1 FOREIGN KEY (user_id) REFERENCES Users(id);