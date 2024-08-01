// Package security provides methods and structs related to user authentication.
package security

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lookeme/short-url/internal/app/domain/user"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/utils"
	"go.uber.org/zap"
)

// Authorization encapsulates the user service and provides methods for user authentication.
type Authorization struct {
	userService *user.UsrService
	Log         *logger.Logger
}

// SecretKey is a constant that holds the secret key used for JWT signing and verification.
// It is used in the following methods:
//   - `BuildJWTString`: Generates a JWT token string with the provided user ID and expiration time.
//     Example usage:
//     ```go
//     tokenString, err := auth.BuildJWTString(userID)
//     ```
//   - `GetUserID`: Parses and verifies a JWT token string and returns the user ID.
//     Example usage:
//     ```go
//     userID := GetUserID(tokenString)
//     ```
//   - `verifyToken`: Verifies the validity of a JWT token string.
//     Example usage:
//     ```go
//     isValid := auth.verifyToken(tokenString)
//     ```
const SecretKey = "secret-key"

// TokenExp is a constant that represents the expiration time for JWT tokens.
// It is used in the following method:
//   - `BuildJWTString`: Generates a JWT token string with the provided user ID and expiration time.
//     Example usage:
//     ```go
//     tokenString, err := auth.BuildJWTString(userID)
//     ```
const TokenExp = time.Hour * 3

// Claims represents the custom claims for JWT authentication.
// It contains the UserID and the RegisteredClaims from the jwt package.
//
// Usage Example:
// var claims Claims
//
//	jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
//	    return []byte(SecretKey), nil
//	})
//
// userID := claims.UserID
type Claims struct {
	UserID int
	jwt.RegisteredClaims
}

// New constructs a new instance of Authorization with a user service and logger.
func New(userService *user.UsrService, logger *logger.Logger) *Authorization {
	return &Authorization{
		userService: userService,
		Log:         logger,
	}
}

// AuthMiddleware is a middleware function that checks for a valid JWT in the Authorization header of the HTTP request.
// If no token is provided, a new user is created and a JWT is generated and sent back in the response.
func (auth *Authorization) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var bearer = "Bearer "
		token := r.Header.Get("Authorization")
		token, err := utils.GetToken(token)
		if err != nil || !auth.verifyToken(token) {
			usr, err := auth.userService.CreateUser()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, err = auth.BuildJWTString(usr.UserID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			bearer += token
			r.Header.Add("Authorization", bearer)
		} else {
			bearer += token
		}
		w.Header().Set("Authorization", bearer)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// BuildJWTString generates a JWT for a specified user ID.
func (auth *Authorization) BuildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserID retrieves a user ID from a given JWT.
func GetUserID(tokenString string) int {
	var claims Claims
	jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	return claims.UserID
}

// verifyToken is a method that takes a token string as input and verifies its validity using JWT.
// It returns true if the token is valid, otherwise false.
// If there is an error during verification, it logs the error and returns false.
func (auth *Authorization) verifyToken(tokenString string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		auth.Log.Log.Error("Error during verifying token", zap.String("error", err.Error()))
		return false
	}
	return token.Valid
}
