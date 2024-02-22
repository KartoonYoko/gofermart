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
	CREATE TABLE "users" (
		"id" SERIAL PRIMARY KEY,
		"login" varchar,
		"password" varchar,
		"loyality_balance_current" integer,
		"loyality_balance_withdrawn" integer
	);
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE "orders" (
		"id" integer PRIMARY KEY,
		"status" varchar,
		"accrual" integer,
		"user_id" integer,

		CONSTRAINT fk_user_id
		FOREIGN KEY (user_id) 
		REFERENCES users (id)
	);
	`)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE "loayality_points_withdrawals" (
		"order_id" integer PRIMARY KEY,
		"user_id" integer,
		"processed_at" TIMESTAMP,
		"sum" integer,

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
