// Package entdb provides helper for creating pgx dabase for ent.
package db

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/XSAM/otelsql"
	"github.com/ernado/example/internal/ent"
	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// Open new connection.
func Open(ctx context.Context, uri string, t *app.Telemetry) (*ent.Client, error) {
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
	db := stdlib.OpenDBFromPool(pool)
	if t != nil {
		options := []otelsql.Option{
			otelsql.WithMeterProvider(t.MeterProvider()),
		}
		if err := otelsql.RegisterDBStatsMetrics(db, options...); err != nil {
			return nil, errors.Wrap(err, "otelsql.RegisterDBStatsMetrics")
		}
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}
