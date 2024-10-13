//go:build api

package api_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
	"github.com/rikeda71/go-gql-sqlc-template/internal"
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/db"
	"github.com/rikeda71/go-gql-sqlc-template/internal/metrics"
	api "github.com/rikeda71/go-gql-sqlc-template/test/api/helper"
)

var (
	sqlcClient *db.Queries
	Pool       *pgxpool.Pool
	Server     *echo.Echo
)

// Execute api tests in the package
func TestMain(m *testing.M) {
	// start postgresql container with docker
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to connect docker: %v", err)
	}
	// set max wait time to 20 seconds because we want to fail fast if we can't connect to the database
	pool.MaxWait = time.Second * 20

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get pwd: %v", err)
	}

	// setup config
	cnf, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("could not start docker resource: %v", err)
	}

	runOptions := &dockertest.RunOptions{
		// sampe compose.yaml settings
		Repository: "postgres",
		// same CloudSQL of Google Cloud settings
		Tag: "16.3-bullseye",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", cnf.DatabaseName),
			fmt.Sprintf("POSTGRES_USER=%s", cnf.DatabaseUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cnf.DatabasePassword),
			"listen_addresses='*'",
		},
		Mounts: []string{},
	}

	resource, err := pool.RunWithOptions(runOptions,
		func(config *docker.HostConfig) {
			// delete container when it stops
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name:              "no",
				MaximumRetryCount: 0,
			}
		},
	)
	if err != nil {
		log.Fatalf("could not start docker resource: %v", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	cnf.DatabaseHost = strings.Split(hostAndPort, ":")[0]
	cnf.DatabasePort, _ = strconv.Atoi(strings.Split(hostAndPort, ":")[1])

	// Retry to connect to the database until available
	if err := pool.Retry(func() error {
		cnf, err := pgxpool.ParseConfig(cnf.DataSource())
		if err != nil {
			return errors.WithMessage(err, "failed to parse postgresql config")
		}
		p, err := pgxpool.NewWithConfig(context.Background(), cnf)
		if err != nil {
			return errors.WithMessage(err, "failed to open postgresql connection")
		}
		// ping postgresql before running queries
		if err := p.Ping(context.Background()); err != nil {
			return errors.WithMessage(err, "failed to ping postgresql")
		}

		// setup global sqlc client
		sqlcClient = db.New(p)
		Pool = p
		return nil
	}); err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// setup db
	/// migration
	migrationPath := path.Join(pwd, "../../", "db", "migrations")
	err = api.ExecuteSQLsFromDir(migrationPath, Pool, "migration")
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	/// seed
	seedPath := path.Join(pwd, "../../", "db", "seed")
	err = api.ExecuteSQLsFromDir(seedPath, Pool, "seed")
	if err != nil {
		log.Fatalf("could not seed database: %v", err)
	}

	// setup app
	/// setup graphql handler
	gqlHandler, err := internal.NewGraphQLHandler(cnf, sqlcClient, metrics.NewClient())
	if err != nil {
		log.Fatalf("could not create graphql handler: %v", err)
	}
	s := internal.NewServer(cnf.Port, *gqlHandler)
	go func() {
		_ = s.Start(false)
	}()
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	Server = s.Server()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %v", err)
	}

	os.Exit(code)
}
