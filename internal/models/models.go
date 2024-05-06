package models

import (
	"fmt"
)

type Request struct {
	URL string `json:"url"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type Response struct {
	Result string `json:"result"`
}

type ShortenData struct {
	ID            int64  `json:"id,omitempty"`
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

func NewShortenData(id int64, originalURL string, shortURL string) *ShortenData {
	return &ShortenData{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func (s *ShortenData) String() string {
	return fmt.Sprintf("{ id=%d, CorrelationID=%s, shortenURL=%s, originURL=%s  }", s.ID, s.CorrelationID, s.ShortURL, s.OriginalURL)
}
