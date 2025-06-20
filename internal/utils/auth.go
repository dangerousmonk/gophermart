package utils

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenLifeTime    = time.Hour * 24
	secretKeySize    = 32
	UserIDHeaderName = "x-user-id"
	AuthCookieName   = "auth"
)

var (
	errExpiredToken  = errors.New("token: has expired")
	errInvalidToken  = errors.New("token: is invalid")
	errInvalidClaims = errors.New("claims: failed to initialize")
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type JWTAuthenticator struct {
	secretKey string
}

func (claims *Claims) Valid() error {
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return errExpiredToken
	}
	return nil
}

type Authenticator interface {
	CreateToken(userID int, duration time.Duration) (string, error)
	ValidateToken(token string) (*Claims, error)
	SetAuth(userID int, w http.ResponseWriter, r *http.Request) error
}

func (auth *JWTAuthenticator) CreateToken(userID int, duration time.Duration) (string, error) {
	claims, err := NewClaims(userID, duration)
	if err != nil {
		slog.Error("CreateToken failed create claims", slog.Any("err", err))
		return "", errInvalidClaims
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.secretKey))
}

func (auth *JWTAuthenticator) ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			slog.Error("ValidateToken jwt method is not valid", slog.Bool("ok", ok))
			return nil, errInvalidToken
		}
		return []byte(auth.secretKey), nil
	}
	_, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	return claims, nil

}

func NewJWTAuthenticator(secretKey string) (Authenticator, error) {
	if len(secretKey) < secretKeySize {
		return nil, errors.New("invalid secretKey len")
	}
	return &JWTAuthenticator{secretKey}, nil
}

func NewClaims(userID int, duration time.Duration) (*Claims, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID: userID,
	}
	return claims, nil
}

func (auth *JWTAuthenticator) SetAuth(userID int, w http.ResponseWriter, r *http.Request) error {
	token, err := auth.CreateToken(userID, tokenLifeTime)
	if err != nil {
		slog.Error("SetAuth failed create token", slog.Any("err", err))
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  AuthCookieName,
		Value: token,
		Path:  "/",
	})
	r.Header.Set(UserIDHeaderName, strconv.Itoa(userID))
	w.Header().Set(UserIDHeaderName, strconv.Itoa(userID))
	return nil
}
