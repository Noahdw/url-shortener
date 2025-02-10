package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"log/slog"
	"net/url"
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
	queries *repository.Queries
	baseURL string
}

func NewURLService(queries *repository.Queries, baseURL string) *urlService {
	return &urlService{
		queries: queries,
		baseURL: baseURL,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, url string, creatorIP string) (string, error) {
	url, err := isValidUrl(url)
	if err != nil {
		slog.Error("Invalid URL provided",
			"error", err,
			"url", url)
		return "", ErrInvalidURL
	}

	short_code := hashString(url)[:6]
	short_url := fmt.Sprintf("http://%s/%s", s.baseURL, short_code)

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

// HandleCreateUrl godoc
// @Summary Redirect to original URL
// @Description Redirects to the original URL associated with the given short code
// @Tags urls
// @Accept json
// @Produce json
// @Param shortCode path string true "Short code of the URL"
// @Success 303 {string} string "Redirect to original URL"
// @Failure 404 {string} string "Short code not found"
// @Router /{shortCode} [get]
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

func isValidUrl(rawUrl string) (string, error) {
	rawUrl = strings.ToLower(rawUrl)
	if !(strings.HasPrefix(rawUrl, "http://") || strings.HasPrefix(rawUrl, "https://")) {
		rawUrl = fmt.Sprintf("https://%s", rawUrl)
	}
	parsed_url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	if parsed_url.Hostname() == "" {
		return "", fmt.Errorf("invalid hostname for %s", rawUrl)

	}
	if !strings.Contains(rawUrl, ".") {
		return "", fmt.Errorf("invalid hostname: missing domain extension for %s", parsed_url.Hostname())
	}
	return rawUrl, nil
}

func hashString(s string) string {
	// Create a new FNV-1a hash
	h := fnv.New64a()

	// Write the string to the hash
	h.Write([]byte(s))

	// Get the hash as bytes
	hashBytes := h.Sum(nil)

	// Convert to base64 string for readability
	// (you could also use hex encoding or other formats)
	return base64.URLEncoding.EncodeToString(hashBytes)
}
