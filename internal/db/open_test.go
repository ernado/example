package db

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func openClient(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}
	cfg.MaxConns = 20
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Minute * 2
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.NewWithConfig")
	}
	return pool, nil
}
