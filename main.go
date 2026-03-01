package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/caioandre182/api-users/api"
	"github.com/caioandre182/api-users/store/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		return
	}

	slog.Info("all system offline")
}

func run() error {
	db := openDB()
	defer db.Close()

	pgStore := postgres.New(db)
	h := api.New(pgStore)

	return http.ListenAndServe(":8080", h.Router())
}

func openDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("missing DATABASE_URL")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	return db
}
