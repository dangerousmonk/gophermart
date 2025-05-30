package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/service/mocks"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newLongPassword() string {
	letter := "a"
	return strings.Repeat(letter, 100)
}

func TestRegisterUserHandler(t *testing.T) {
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
					CreateUser(gomock.Any(), gomock.AssignableToTypeOf(&models.UserRequest{Login: "fake_login@example.com", HashedPassword: gomock.Any().String(), Password: "super_secret"})).
					Times(1).
					Return(1, nil)
			},
		},
		{
			name:         "Validation errors",
			body:         models.UserRequest{Login: "", Password: ""},
			expectedCode: http.StatusBadRequest,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Password too long",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: newLongPassword()},
			expectedCode: http.StatusInternalServerError,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "User already registered",
			body:         models.UserRequest{Login: "fake_login@example.com", Password: "super_secret"},
			expectedCode: http.StatusConflict,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.AssignableToTypeOf(&models.UserRequest{Login: "fake_login@example.com", HashedPassword: gomock.Any().String(), Password: "super_secret"})).
					Times(1).
					Return(0, service.ErrLoginExists)
				r.EXPECT().IsUniqueViolation(gomock.Any(), gomock.Any()).Times(1)
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

			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(json))
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			s := service.NewGophermartService(repo, &cfg)
			handler := NewHandler(*s, jwtAuthenticator)

			handler.RegisterUser(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}

		})
	}

}
