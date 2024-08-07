// Package user defines structures and methods related to user management.
package user

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/storage"
	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// UsrService represents a service for user management.
// It incorporates logging and a repository for user persistence.
type UsrService struct {
	userRepository storage.UserRepository
	Log            *logger.Logger
}

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"

// NewUserService constructs a new instance of UsrService with necessary dependencies.
// It requires a user repository for data persistence and a logger.
func NewUserService(userRepository storage.UserRepository, log *logger.Logger) UsrService {
	return UsrService{
		userRepository,
		log,
	}
}

// CreateUser generates a new user with a random username and password,
// hashes the password using Argon2 algorithm, and saves the user into the repository.
// It returns a User model or an error if the operation fails.
func (s *UsrService) CreateUser() (models.User, error) {
	user := models.User{}
	strPass := generatePass(8)
	strName := generatePass(5)
	hash, err := generateFromPassword(strPass)
	if err != nil {
		return models.User{}, err
	}
	user.Name = strName
	ID, err := s.userRepository.SaveUser(strName, hash)
	if err != nil {
		return models.User{}, err
	}
	user.UserID = ID
	return user, err
}

// FindByID retrieves a user from the repository by their ID.
func (s *UsrService) FindByID(userID int) (models.User, error) {
	return s.userRepository.FindByID(userID)
}

// generateFromPassword creates a password hash using the Argon2 ID hashing algorithm.
// It returns a base64-encoded hash.
func generateFromPassword(password string) (b64Hash string, err error) {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	b64Hash = base64.RawStdEncoding.EncodeToString(hash)
	return b64Hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generatePass(length int) string {
	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}
