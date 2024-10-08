package utils

// ShortToken is string of the number of BASE64 symbols from the url safe alphabet.
// It represent 6*len(token) bits of data. ShortToken is not correct BASE64 data
// representation as number of bits is not always a multiple of 8 (1 byte).

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// ShortToken - interface for short token creation
type ShortToken interface {
	Get() string        // returns new random short token
	Check(string) error // check the token length and alphabet
}

// ShortToken - interface for short token creation
type shortToken struct {
	length  int // token length
	bufSize int // bytes bufer size
}

// NewShortToken returns new ShortToken instance
func NewShortToken(length int) ShortToken {
	return &shortToken{
		length:  length,
		bufSize: length*6/8 + 1,
	}
}

// Get creates the token from random or debugging source
func (s *shortToken) Get() string {

	// prepare bytes bufer
	buf := make([]byte, s.bufSize)

	// get secure random bytes
	n, err := rand.Read(buf)
	if err != nil || n != s.bufSize {
		panic(fmt.Errorf("error while retriving random data: %d %v", n, err.Error()))
	}
	// return shortened to tokenLenS BASE64 representation
	return base64.URLEncoding.EncodeToString(buf)[:s.length]
}

// Check checks the lenght of token and its alphabet
func (s *shortToken) Check(sToken string) error {

	// check length
	if len(sToken) != s.length {
		return errors.New("wrong token length")
	}

	// check base64 alphabet
	if _, err := base64.URLEncoding.DecodeString(sToken + "AAAA"[:4-s.length%4]); err != nil {
		return errors.New("wrong token alphabet")
	}
	return nil
}

// GetToken - function to create token
func GetToken(str string) (string, error) {
	if str == "" {
		return "", errors.New("token is invalid")
	}
	tokens := strings.Split(str, " ")
	if len(tokens) != 2 {
		return "", errors.New("token is invalid")
	}
	return tokens[1], nil
}

// CreateShortURL - function to create short url
func CreateShortURL(key string, baseURL string) string {
	return fmt.Sprintf("%s/%s", baseURL, key)
}

// ErrorCode - compare errors
func ErrorCode(err error) string {
	var pgerr *pgconn.PgError
	ok := errors.As(err, &pgerr)
	if !ok {
		return ""
	}
	return pgerr.Code
}
