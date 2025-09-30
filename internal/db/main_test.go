package db

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBTestSuite struct {
	suite.Suite
	db   *DB
	pgx  *pgx.Conn
	pool *pgxpool.Pool
	uri  string
}

func (suite *DBTestSuite) SetupSuite() {
	// Skip if running on non-linux architecture.
	ctx := suite.T().Context()
	if runtime.GOOS != "linux" {
		// TODO: support macos with colima.
		suite.T().Skip("Skipping test on non-linux architecture")
		return
	}
	const (
		dbName     = "test_db"
		dbUser     = "test_user"
		dbPassword = "test_password"
	)
	var env = map[string]string{
		"POSTGRES_PASSWORD": dbPassword,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       dbName,
	}
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env:          env,
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	suite.Require().NoError(err)

	ip, err := container.ContainerIP(ctx)
	suite.Require().NoError(err)

	//nolint:nosprintfhostport
	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, ip, dbName,
	)

	suite.T().Logf("Postgres URI: %s", uri)

	pool, err := openClient(ctx, uri)
	suite.Require().NoError(err)

	//nolint:nosprintfhostport
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, ip, "postgres",
	))
	suite.Require().NoError(err)

	suite.db = New(pool)
	suite.uri = uri
	suite.pool = pool
	suite.pgx = conn
}

func (suite *DBTestSuite) migrate() *migrate.Migrate {
	//nolint:nosprintfhostport
	migrateURI := strings.ReplaceAll(suite.uri, "postgres://", "pgx5://")
	sourceURI := "file://_migrations"

	m, err := migrate.New(sourceURI, migrateURI)
	suite.Require().NoError(err)
	return m
}

func (suite *DBTestSuite) closeMigrate(m *migrate.Migrate) {
	if m == nil {
		return
	}
	sourceErr, dbErr := m.Close()
	suite.Require().NoError(sourceErr)
	suite.Require().NoError(dbErr)
}

func (suite *DBTestSuite) SetupTest() {
	m := suite.migrate()
	suite.Require().NoError(m.Up())
	suite.closeMigrate(m)
}

func (suite *DBTestSuite) TearDownTest() {
	if suite.pool == nil {
		return
	}

	suite.pool.Close()

	// Drop and recreate the database for each test.
	ctx := suite.T().Context()
	_, err := suite.pgx.Exec(ctx, "DROP DATABASE IF EXISTS test_db")
	suite.Require().NoError(err)
	_, err = suite.pgx.Exec(ctx, "CREATE DATABASE test_db WITH OWNER test_user")
	suite.Require().NoError(err)

	pool, err := openClient(ctx, suite.uri)
	suite.Require().NoError(err)

	suite.pool = pool
	suite.db = New(suite.pool)
}

func TestDBTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DBTestSuite))
}
