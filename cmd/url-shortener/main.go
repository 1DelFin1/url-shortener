package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	logger.Info("info message", slog.String("env", cfg.Env))
	logger.Debug("debug message")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://yandex.ru", "yandex")
	if err != nil {
		logger.Error("Failed to save url", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("saved: ", slog.Int64("id", id))

	id2, err := storage.SaveURL("https://goole.com", "google")
	if err != nil {
		logger.Error("Failed to save url", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("saved: ", slog.Int64("id", id2))

	_ = storage

	// TODO: init router: chi

	// TODO: run server
}
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return logger
}
