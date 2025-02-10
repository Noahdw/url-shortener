package main

// @title URL Shortener API
// @version 1.0
// @description A service to create and manage shortened URLs
// @host localhost:8080
// @BasePath /

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Noahdw/url-shortener/cmd/api/docs"
	"github.com/Noahdw/url-shortener/internal/app"
	httphandler "github.com/Noahdw/url-shortener/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger) // Set as default logger

	dbConfig := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		slog.Error("Unable to create connection pool",
			"error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	app := app.NewApp(dbpool)

	urlHandler := httphandler.NewURLHandler(app.URLService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/generateurl", urlHandler.HandleGenerateShortCode)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	r.Get("/*", urlHandler.HandleUrlRedirect)

	slog.Info("Starting server", "port", ":8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
