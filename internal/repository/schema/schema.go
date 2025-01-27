package schema

const (
	DBSchema = `
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

		CREATE INDEX IF NOT EXISTS accounts_user_id ON orders (user_id);

		ALTER TABLE Orders ADD CONSTRAINT Orders_fk0 FOREIGN KEY (user_id) REFERENCES Users(id);
		ALTER TABLE Accounts ADD CONSTRAINT Accounts_fk1 FOREIGN KEY (user_id) REFERENCES Users(id);
	`

	SelectUser = `
		SELECT * 
		FROM users 
		WHERE login = @login;
	`

	InsertUser = `
		INSERT INTO users (login, password)
		VALUES (@login, @password)
		RETURNING id;
	`

	SelectOrders = `
		SELECT 
			id as number,
			user_id,
			status,
			accrual::money::numeric::float8,
			sum::money::numeric::float8,
			date
		FROM orders
	 	WHERE user_id = @user_id AND sum::money::numeric = 0 AND status IN ('NEW', 'PROCESSING', 'PROCESSED' ,'INVALID')
		ORDER BY date DESC
		FOR UPDATE;
	`

	InsertOrder = `
		INSERT INTO orders (id, user_id, status, date)
		VALUES (@id, @user_id, @status, @date);
	`

	SelectOrderFromAnotherUser = `
		SELECT * 
		FROM orders 
		WHERE id = @id AND user_id <> @user_id
		FOR UPDATE;
	`

	SelectUserBalanceWithdrawn = `
		SELECT 
			a.amount::money::numeric as current,
			sum(b.sum)::money::numeric as withdrawn 
		FROM accounts a LEFT JOIN orders b USING (user_id)
		WHERE a.user_id = @user_id AND b.user_id = @user_id
		GROUP BY a.amount;
	`

	SelectUserBalance = `
		SELECT amount::money::numeric as current
		FROM accounts
		WHERE user_id = @user_id
		FOR UPDATE;
	`

	InsertOrderWithdraw = `
		INSERT INTO orders (id, user_id, sum,  date)
		VALUES (@id, @user_id, @sum, @date);
	`

	UpdateBalanceWithdraw = `
		UPDATE accounts
		SET amount = accounts.amount - @sum
		WHERE user_id = @user_id;
	`

	SelectWithdrawals = `
		SELECT
			id as order,
			sum::money::numeric,
			date as processed_at
		FROM orders
		WHERE user_id = @user_id AND sum::money::numeric > 0
		ORDER BY date DESC
		FOR UPDATE;
	`

	UpdateOrderStatus = `
		UPDATE orders
		SET status = @status
		WHERE user_id = @user_id;
	`

	UpsertBalanceAccrual = `
		INSERT INTO accounts (user_id, amount)
		VALUES (@user_id, @accrual)
		ON CONFLICT(user_id)  
		DO UPDATE
		SET amount = accounts.amount + @accrual;
	`

	UpdateOrderAccrual = `
		UPDATE orders
		SET
			status = @status,
			accrual = orders.accrual + @accrual
		WHERE user_id = @user_id;
	`

	SelectUnprocessedOrders = `
		SELECT
			id as number,
			user_id,
			status,
			accrual::money::numeric::float8,
			sum::money::numeric::float8,
			date
		FROM orders
		WHERE status IN ('NEW', 'PROCESSING', 'REGISTERED')
		FOR UPDATE;
	`
)
