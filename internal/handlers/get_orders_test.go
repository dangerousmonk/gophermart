package handlers

import (
	"context"
	"errors"
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

func TestGetOrdersHandler(t *testing.T) {
	cfg := config.Config{ServerAddr: "http://localhost:8080", JWTSecret: "foobarfoobarfoobarfoobarfoobafoobarfoobarfoobar"}
	jwtAuthenticator, err := utils.NewJWTAuthenticator(cfg.JWTSecret)
	require.NoError(t, err)

	orders := make([]models.Order, 0)
	orders = append(orders, newFakeOrder(1, "4111111111111111", models.StatusNew))

	testCases := []struct {
		name         string
		expectedCode int
		buildStubs   func(s *mocks.MockRepository)
		userID       int
	}{
		{
			name:         "Happy case",
			expectedCode: http.StatusOK,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUserOrders(gomock.Any(), 1).
					Times(1).
					Return(orders, nil)
			},
		},
		{
			name:         "Some error",
			expectedCode: http.StatusInternalServerError,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUserOrders(gomock.Any(), 1).
					Times(1).
					Return(nil, errors.New("some error"))
			},
		},
		{
			name:         "No orders",
			expectedCode: http.StatusNoContent,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetUserOrders(gomock.Any(), 1).
					Times(1).
					Return(make([]models.Order, 0), nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			s := service.NewGophermartService(repo, &cfg)

			handler := NewHandler(*s, jwtAuthenticator)
			handler.GetUserOrders(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}

		})
	}

}
