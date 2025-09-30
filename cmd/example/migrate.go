package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ernado/example/internal/db"
	"github.com/go-faster/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
)

func Migrate() *cobra.Command {
	return &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := iofs.New(db.Migrations, "_migrations")
			if err != nil {
				return errors.Wrap(err, "create iofs driver")
			}
			uri := strings.ReplaceAll(os.Getenv("DATABASE_URL"), "postgres://", "pgx5://")
			if uri == "" {
				uri = "pgx5://postgres:postgres@localhost:5432/example?sslmode=disable"
			}
			m, err := migrate.NewWithSourceInstance("iofs", d, uri)
			if err != nil {
				return errors.Wrap(err, "create migrate")
			}
			if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return errors.Wrap(err, "migrate up")
			}
			sourceErr, dbErr := m.Close()
			if sourceErr != nil {
				return errors.Wrap(sourceErr, "close source")
			}
			if dbErr != nil {
				return errors.Wrap(dbErr, "close db")
			}

			fmt.Println("Migration completed successfully")

			return nil
		},
	}
}
