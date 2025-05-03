package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDHeaderName string     = "x-user-id"
	UserIDContextKey contextKey = "userID"
	AuthCookieName              = "auth"
)

func AuthMiddleware(authenticator utils.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(AuthCookieName)

			if err != nil {
				slog.Error("AuthMiddleware error", slog.Any("err", err))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := getUserFromCookie(authenticator, r, cookie)
			if err != nil {
				slog.Error("AuthMiddleware failed resolve cookie", slog.Any("err", err))

				if errors.Is(err, jwt.ErrTokenExpired) {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctxUser := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctxUser))
		}
		return http.HandlerFunc(fn)
	}
}

func getUserFromCookie(auth utils.Authenticator, r *http.Request, cookie *http.Cookie) (int, error) {
	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		return 0, err
	}
	userID := claims.UserID
	r.Header.Set(UserIDHeaderName, strconv.Itoa(userID))
	return userID, nil
}
