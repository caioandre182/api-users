package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/caioandre182/api-users/api"
	"github.com/caioandre182/api-users/domain"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		return
	}

	slog.Info("all system offline")
}

func run() error {
	db := make(map[string]domain.User)
	handler := api.NewRouter(db)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	db["1"] = domain.User{ID: "1", FirstName: "Caio", LastName: "Macedo", Biography: "teste"}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
