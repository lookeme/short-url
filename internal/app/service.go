// Package app provides the main interfaces for the application.
package app

import "github.com/lookeme/short-url/internal/models"

// ShortenURLService provides an interface to operate on ShortenData.
type ShortenURLService interface {

	// CreateAndSave creates a shortened URL and stores it,
	// it requires a key and the user ID as inputs and returns the created URL.
	CreateAndSave(key string, userID int) (string, error)

	// FindByURL searches for an existing ShortenData entry using the given key.
	// It returns the corresponding ShortenData and a boolean indicating if the entry exists.
	FindByURL(key string) (models.ShortenData, bool)

	// FindByKey searches for an existing ShortenData entry using the provided key.
	// It returns the corresponding ShortenData and a boolean indicating if the entry exists.
	FindByKey(key string) (models.ShortenData, bool)

	// FindAll returns all available ShortenData within the database.
	FindAll() ([]models.ShortenData, error)

	// CreateAndSaveBatch creates a batch of shorten URLs and saves them,
	// requires an array of BatchRequest as input and returns an array of BatchResponse.
	CreateAndSaveBatch(urls []models.BatchRequest) ([]models.BatchResponse, error)

	// DeleteByShortURLs deletes the ShortenData entries whose keys are in the given URLs.
	// Returns an error if it fails.
	DeleteByShortURLs(urls []string) error
}

// UserService provides an interface for operations on User models.
type UserService interface {

	// CreateUser is a function that creates a new user given a username,
	// and returns the newly created User model.
	CreateUser(userName string) (models.User, error)

	// FindByID searches for an existing User entry given a user ID,
	// and returns the corresponding User model and a boolean indicating if the entry exists.
	FindByID(userID int) (models.User, bool)
}
