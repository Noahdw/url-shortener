package app

import (
	"github.com/Noahdw/url-shortener/internal/repository"
	"github.com/Noahdw/url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool       *pgxpool.Pool
	URLService service.URLService
}

func NewApp(pool *pgxpool.Pool) *App {
	repo := repository.New(pool)
	urlService := service.NewURLService(repo, "localhost:8080")
	return &App{
		pool:       pool,
		URLService: urlService,
	}
}
