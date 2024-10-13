package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rikeda71/go-gql-sqlc-template/internal"
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/db"
	"github.com/rikeda71/go-gql-sqlc-template/internal/metrics"
)

const (
	LogCountTotal = "log_count"
	level         = "level"
)

func main() {
	cnf, err := internal.NewConfig()
	if err != nil {
		panic(err)
	}

	// metrics
	m := metrics.NewClient()
	m.RegisterCounter(LogCountTotal, "ログの出現回数", level)

	// logger
	logLevel := slog.LevelInfo
	if cnf.DebugMode {
		logLevel = slog.LevelDebug
	}
	opt := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// count by log level
			switch {
			case a.Key == slog.LevelKey:
				m.Count(LogCountTotal, 1, a.Value.String())
			}
			return a
		},
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, opt)))

	// infrastructure
	/// db
	poolCnf, err := pgxpool.ParseConfig(cnf.DataSource())
	if err != nil {
		slog.Error("failed to connect db", "error", err.Error())
		panic(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCnf)
	if err != nil {
		slog.Error("failed to connect db", "error", err.Error())
		panic(err)
	}
	defer func() {
		pool.Close()
	}()
	q := db.New(pool)

	// presentation
	gqlHandler, err := internal.NewGraphQLHandler(cnf, q, m)
	if err != nil {
		panic(err)
	}
	s := internal.NewServer(cnf.Port, *gqlHandler)
	go func() {
		if err := s.Start(cnf.DebugMode); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("could not start server.", "err", err.Error())
		}
	}()

	// graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func(server internal.Server) {
		// wait for signal
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		/// block until signal received
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.GracefulTimeout)*time.Second)
		defer cancel()

		slog.Info("shutting down server with graceful...")
		defer wg.Done()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed.", "err", err.Error())
		}
	}(*s)
	wg.Wait()

	slog.Info("server shutdown.")
}
