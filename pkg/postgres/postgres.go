package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL string
}

func NewPool(config *Config) (*pgxpool.Pool, error) {
	pgpool, err := pgxpool.New(context.Background(), config.URL)
	if err != nil {
		return nil, err
	}
	return pgpool, nil
}
