package service

import (
	"context"
	"testing"

	"github.com/Noahdw/url-shortener/internal/repository"
)

func TestValidURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"", false},
		{".com", false},
		{".", false},
		{"www.", false},
		{"http://www.com", false},
		{"http://reddit", false},
		{"reddit", false},
		{"http://reddit.com", true},
		{"https://reddit.com", true},
		{"reddit.com", true},
		{"url.net", true},
		{"blog.google.com", true},
	}

	for _, test := range tests {
		_, err := validatedURL(test.url)
		actual := true
		if err != nil {
			actual = false
		}

		if test.expected != actual {
			t.Errorf("url %q valid marked as %t, expected %t", test.url, actual, test.expected)
		}
	}
}

func TestValidatedURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"reddit.com", "https://www.reddit.com"},
		{"www.reddit.com", "https://www.reddit.com"},
		{"https://www.reddit.com", "https://www.reddit.com"},
	}

	for _, test := range tests {
		actual, _ := validatedURL(test.url)

		if test.expected != actual {
			t.Errorf("validated url for base %q: actual %q, expected %q", test.url, actual, test.expected)
		}
	}
}

func TestCreateShortURL(t *testing.T) {
	baseURLS := []string{
		"localhost:8080",
		"https://www.surl.com",
	}
	for _, baseURL := range baseURLS {
		tests := []struct {
			originalURL string
			expectedURL string
		}{
			{"", ""},
			{"reddit.com", baseURL + "/AZXhKF"},
			{"https://www.reddit.com", baseURL + "/AZXhKF"},
			{"reddit", ""},
			{"www.wsj.com", baseURL + "/Jk_M4s"},
		}

		db := repository.NewRepoMock()
		for _, test := range tests {
			service := NewURLService(db, baseURL)
			actualURL, _ := service.CreateShortURL(context.Background(), test.originalURL, "")
			if actualURL != test.expectedURL {
				t.Errorf("Creating short code for URL %q: actual %q, expected %q", test.originalURL, actualURL, test.expectedURL)
			}
		}
	}
}
