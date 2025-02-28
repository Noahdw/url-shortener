package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"log/slog"
	"strings"
	"time"

	"github.com/Noahdw/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type URLService interface {
	CreateShortURL(ctx context.Context, url string, creatorIP string) (string, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}

type urlService struct {
	queries repository.Querier
	baseURL string
}

func NewURLService(queries repository.Querier, baseURL string) *urlService {
	return &urlService{
		queries: queries,
		baseURL: baseURL,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, url string, creatorIP string) (string, error) {
	url, err := validatedURL(url)
	if err != nil {
		slog.Error("Invalid URL provided",
			"error", err,
			"url", url)
		return "", ErrInvalidURL
	}

	short_code := hashString(url)[:6]

	short_url := fmt.Sprintf("%s/%s", s.baseURL, short_code)

	_, err = s.queries.CreateUrlMapping(ctx, repository.CreateUrlMappingParams{
		OriginalUrl: url,
		ShortCode:   short_code,
		ExpiresAt: pgtype.Timestamptz{
			Time:             time.Now().Add(time.Minute),
			InfinityModifier: pgtype.Finite,
			Valid:            true},
		CreatorIp: pgtype.Text{String: creatorIP, Valid: true},
	})
	if err != nil {
		// 23505 = unique key violation - so get existing key
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			slog.Info("Fetching existing short code for URL",
				"url", url)
			short_code, err = s.queries.GetShortCodeFromOriginalUrl(ctx, url)
			if err != nil {
				return "", ErrNotFound
			}
		} else {
			return "", fmt.Errorf("failed to create URL mapping: %w", err)
		}
	}

	slog.Info("Generated short url", url, short_code)
	return short_url, nil
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.queries.GetOriginalUrlFromShortCode(ctx, shortCode)
	if err != nil {
		return "", ErrNotFound
	}
	slog.Info("Redirect to URL",
		"short_code", shortCode,
		"url", url)
	return url, nil
}

func validatedURL(rawUrl string) (string, error) {
	if len(rawUrl) > 256 { // length we store in db
		return "", fmt.Errorf("URL too long")
	}

	var scheme string
	if strings.HasPrefix(rawUrl, "http://") {
		scheme = "http://"
	} else if strings.HasPrefix(rawUrl, "https://") {
		scheme = "https://"
	}
	if len(scheme) > 0 {
		rawUrl = rawUrl[len(scheme):]
	} else {
		scheme = "https://"
	}

	rawUrl = strings.TrimPrefix(rawUrl, "www.")

	if strings.Count(rawUrl, ".") == 0 {
		return "", fmt.Errorf("no domain dot")
	}
	if string(rawUrl[0]) == "." {
		return "", fmt.Errorf("bad domain dot")
	}
	return scheme + "www." + rawUrl, nil
}

func hashString(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	hashBytes := h.Sum(nil)

	return base64.URLEncoding.EncodeToString(hashBytes)
}
