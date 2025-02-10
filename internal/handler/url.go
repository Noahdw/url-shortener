package httphandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Noahdw/url-shortener/internal/service"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}

}

// HandleGenerateShortCode godoc
// @Summary Generate short URL
// @Description Creates a shortened URL from a provided original URL
// @Tags urls
// @Accept json
// @Produce json
// @Param url query string true "Original URL to shorten"
// @Success 200 {string} string "Short URL"
// @Failure 400 {string} string "Invalid URL format"
// @Router /generateurl [get]
func (h *URLHandler) HandleGenerateShortCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()
	url := params.Get("url")
	ip := r.RemoteAddr
	short_url, err := h.service.CreateShortURL(ctx, url, ip)
	if err != nil {
		h.handleError(w, r, err)
	} else {
		io.WriteString(w, short_url)
	}
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
func (h *URLHandler) HandleUrlRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	short_code := r.URL.Path[1:]
	url, err := h.service.GetOriginalURL(ctx, short_code)
	if err != nil {
		h.handleError(w, r, err)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (h *URLHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, service.ErrInvalidURL):
		status = http.StatusBadRequest
		message = "Invalid URL format"
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
		message = "Short code not found"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
		// Log the full error internally
		slog.Error("unexpected error",
			"error", err,
			"path", r.URL.Path,
			"method", r.Method)
	}

	http.Error(w, message, status)
}
