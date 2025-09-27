package entdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/ernado/example/internal/ent"
)

type DBTestSuite struct {
	suite.Suite
	ent  *ent.Client
	db   *DB
	pgx  *pgx.Conn
	pool *pgxpool.Pool
	uri  string
}

func (suite *DBTestSuite) SetupSuite() {
	ctx := suite.T().Context()

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

	client, pool, err := openClient(ctx, uri)
	suite.Require().NoError(err)

	//nolint:nosprintfhostport
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, ip, "postgres",
	))
	suite.Require().NoError(err)

	suite.pgx = conn
	suite.ent = client
	suite.db = New(suite.ent)
	suite.uri = uri
	suite.pool = pool
}

func (suite *DBTestSuite) TearDownTest() {
	if suite.ent == nil {
		return
	}

	err := suite.ent.Close()
	suite.Require().NoError(err)

	suite.pool.Close()

	// Drop and recreate the database for each test.
	ctx := suite.T().Context()
	_, err = suite.pgx.Exec(ctx, "DROP DATABASE IF EXISTS test_db")
	suite.Require().NoError(err)
	_, err = suite.pgx.Exec(ctx, "CREATE DATABASE test_db WITH OWNER test_user")
	suite.Require().NoError(err)

	// Migrate.
	client, pool, err := openClient(ctx, suite.uri)
	suite.Require().NoError(err)

	suite.ent = client
	suite.pool = pool
	suite.db = New(suite.ent)

	err = suite.ent.Schema.Create(ctx)
	suite.Require().NoError(err)
}

func TestDBTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DBTestSuite))
}
