package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/KarmaBeLike/crypto-service/config"
	"github.com/KarmaBeLike/crypto-service/internal/database"
	"github.com/KarmaBeLike/crypto-service/internal/handlers"
	"github.com/KarmaBeLike/crypto-service/internal/repository"
	"github.com/KarmaBeLike/crypto-service/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
		return
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("error running migrations", slog.Any("error", err))
		return
	}

	tokenRepo := repository.NewTokenRepository(db)
	tokenService := service.NewTokenService(tokenRepo)
	tokenHandler := handlers.NewTokenHandler(tokenService)

	r := mux.NewRouter()
	r.HandleFunc("/tokens", tokenHandler.GetAndStoreTokens).Methods(http.MethodGet)

	// Запускаем сервер
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
