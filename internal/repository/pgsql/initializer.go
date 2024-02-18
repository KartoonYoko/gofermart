package pgsql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type initializer struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) (*initializer, error) {
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

	tx.ExecContext(ctx, `
	CREATE TABLE "users" (
		"id" integer PRIMARY KEY,
		"login" varchar,
		"password" varchar,
		"loyality_balance_current" integer,
		"loyality_balance_withdrawn" integer
	);
	`)

	tx.ExecContext(ctx, `
	CREATE TABLE "orders" (
		"id" long PRIMARY KEY,
		"status" varchar,
		"accrual" integer,
		"user_id" integer,

		CONSTRAINT fk_user_id
		FOREIGN KEY (user_id) 
		REFERENCES users (id),
	);
	`)

	tx.ExecContext(ctx, `
	CREATE TABLE "loayality_points_withdrawals" (
		"order_id" PRIMARY KEY,
		"user_id" integer,
		"processed_at" datetime,
		"sum" integer,

		CONSTRAINT fk_user_id
		FOREIGN KEY (user_id) 
		REFERENCES users (id),
	);
	`)

	return tx.Commit()
}
