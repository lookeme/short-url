package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type ShortenData struct {
	ID          uuid.UUID `json:"uuid,omitempty"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func NewShortenData(originalURL string, shortURL string) *ShortenData {
	id := uuid.New()
	return &ShortenData{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func (s *ShortenData) String() string {
	return fmt.Sprintf("{ id=%s, shortenURL=%s, originURL=%s  }", s.ID, s.ShortURL, s.OriginalURL)
}
