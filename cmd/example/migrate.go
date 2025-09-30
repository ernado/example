package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func Migrate() *cobra.Command {
	return &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := os.Getenv("DATABASE_URL")
			uri = strings.ReplaceAll(uri, "postgres://", "pgx5://")
			if uri == "" {
				uri = "pgx5://postgres:postgres@localhost:5432/example?sslmode=disable"
			}
			sourceURI := "file://internal/dbraw/_migrations"
			m, err := migrate.New(sourceURI, uri)
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
