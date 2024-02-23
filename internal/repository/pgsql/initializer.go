package pgsql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type initializer struct {
	conn *sqlx.DB
}

func NewInitializer(ctx context.Context, db *sqlx.DB) (*initializer, error) {
	i := &initializer{
		conn: db,
	}

	return i, nil
}

func (in *initializer) Init(ctx context.Context) error {
	tx, err := in.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// вознаграждение начисляется и тратится в виртуальных баллах из расчёта 1 балл = 1 единица местной валюты.
	// TODO посмотреть monetary types для posgresql

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS "users" (
		"id" SERIAL PRIMARY KEY,
		"login" VARCHAR,
		"password" VARCHAR,
		"loyality_balance_current" INTEGER,
		"loyality_balance_withdrawn" INTEGER
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login ON users(login);
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS "orders" (
		"id" INTEGER PRIMARY KEY,
		"status" VARCHAR,
		"accrual" INTEGER,
		"user_id" INTEGER,

		CONSTRAINT fk_user_id
		FOREIGN KEY (user_id) 
		REFERENCES users (id)
	);
	`)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS "loayality_points_withdrawals" (
		"order_id" INTEGER PRIMARY KEY,
		"user_id" INTEGER,
		"processed_at" TIMESTAMP,
		"sum" INTEGER,

		CONSTRAINT fk_user_id
		FOREIGN KEY (user_id) 
		REFERENCES users (id)
	);
	`)

	if err != nil {
		return err
	}

	return tx.Commit()
}
