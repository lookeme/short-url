// Package models defines the data structures used for managing URL shortening and user details.
package models

import (
	"fmt"
)

// Request represents the request for URL shortening with a single URL.
type Request struct {
	URL string `json:"url"`
}

// BatchRequest represents a request for URL shortening with multiple URLs, each associated with a correlation ID.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse represents the response to a BatchRequest, including the correlation ID and the shortened URL.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Response represents a generally applicable response with a result string.
type Response struct {
	Result string `json:"result"`
}

// ShortenData encapsulates all data associated with a shortened URL, including metadata like the user ID who created it.
type ShortenData struct {
	ID            int64  `json:"-"`
	CorrelationID string `json:"-"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	UserID        int    `json:"-"`
	DeletedFlag   bool   `db:"is_deleted"`
}

// User represents the details of a user in the system, including their ID, name, and password.
type User struct {
	UserID   int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Pass     string `json:"pass,omitempty"`
	IsActive bool   `json:"is_active"`
}

// NewShortenData creates a new instance of ShortenData with the given parameters.
func NewShortenData(id int64, originalURL string, shortURL string, userID int) *ShortenData {
	return &ShortenData{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
}

// String provides a string representation of the ShortenData instance.
func (s *ShortenData) String() string {
	return fmt.Sprintf("{ id=%d, CorrelationID=%s, shortenURL=%s, originURL=%s  }", s.ID, s.CorrelationID, s.ShortURL, s.OriginalURL)
}
