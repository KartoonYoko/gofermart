package pgsql

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewSQLxConnection создаёт экземпляр sqlx.DB и накатывает миграции
func NewSQLxConnection(ctx context.Context, cfg *Config) (*sqlx.DB, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()
	db, err := sqlx.Connect("pgx", cfg.connectionString())
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
