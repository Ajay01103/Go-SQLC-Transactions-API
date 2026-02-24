package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Ajay01103/goTransactonsAPI/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=123 dbname=godb sslmode=disable"),
		},
		jwtSecret: env.GetString("JWT_SECRET", "change-me-in-production"),
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}