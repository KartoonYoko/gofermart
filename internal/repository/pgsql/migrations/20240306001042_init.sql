-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users" (
	"id" SERIAL PRIMARY KEY,
	"login" VARCHAR,
	"password" VARCHAR,
	"loyality_balance_current" INTEGER,
	"loyality_balance_withdrawn" INTEGER
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login ON users(login);

CREATE TABLE IF NOT EXISTS "orders" (
	"id" BIGINT PRIMARY KEY,
	"status" VARCHAR,
	"accrual" INTEGER,
	"user_id" INTEGER,
	"created_at" TIMESTAMP without time zone default (now() at time zone 'utc'),

	CONSTRAINT fk_user_id
	FOREIGN KEY (user_id) 
	REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS "withdrawals" (
	"order_id" BIGINT PRIMARY KEY,
	"user_id" INTEGER,
	"processed_at" TIMESTAMP without time zone default (now() at time zone 'utc'),
	"sum" INTEGER,

	CONSTRAINT fk_user_id
	FOREIGN KEY (user_id) 
	REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
