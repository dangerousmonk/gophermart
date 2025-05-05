package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/repository/mocks"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newFakeHashed(password string) string {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return ""
	}
	return hashed
}

func newFakeUser(password string) models.User {
	u := models.User{ID: 1, Login: "fake_login@example.com", Password: newFakeHashed(password)}
	return u
}

func TestLoginUserHandler(t *testing.T) {
	cfg := config.Config{ServerAddr: "http://localhost:8080", JWTSecret: "foobarfoobarfoobarfoobarfoobafoobarfoobarfoobar"}
	jwtAuthenticator, err := utils.NewJWTAuthenticator(cfg.JWTSecret)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		body         models.UserRequest
		expectedCode int
		buildStubs   func(s *mocks.MockRepository)
		userID       int
	}{
		{
			name:         "Happy case",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: "super_secret"},
			expectedCode: http.StatusOK,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUser(gomock.Any(), "fake_login@example.com").
					Times(1).
					Return(newFakeUser("super_secret"), nil)
			},
		},
		{
			name:         "Validation errros",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: "123"},
			expectedCode: http.StatusBadRequest,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "User not found",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: "super_secret"},
			expectedCode: http.StatusUnauthorized,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUser(gomock.Any(), "fake_login@example.com").
					Times(1).
					Return(models.User{}, service.ErrNoUserFound)
				r.EXPECT().IsNoRows(gomock.Any()).Times(1)
			},
		},
		{
			name:         "Wrong password",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: "foo_bar"},
			expectedCode: http.StatusUnauthorized,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUser(gomock.Any(), "fake_login@example.com").
					Times(1).
					Return(newFakeUser("super_secret"), nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			json, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(json))
			_ = context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			w := httptest.NewRecorder()

			s := service.NewGophermartService(repo, &cfg)

			handler := NewHandler(*s, jwtAuthenticator)
			handler.LoginUser(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}

		})
	}

}
