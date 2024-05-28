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

type Authorization struct {
	userService *user.UsrService
	Log         *logger.Logger
}

const SecretKey = "secret-key"
const TokenExp = time.Hour * 3

type Claims struct {
	UserID int
	jwt.RegisteredClaims
}

func New(userService *user.UsrService, logger *logger.Logger) *Authorization {
	return &Authorization{
		userService: userService,
		Log:         logger,
	}
}

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

func GetUserID(tokenString string) int {
	var claims Claims
	jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	return claims.UserID
}

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
