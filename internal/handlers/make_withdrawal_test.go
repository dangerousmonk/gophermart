package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestMakeWithdrawalHandler(t *testing.T) {
	cfg := config.Config{ServerAddr: "http://localhost:8080", JWTSecret: "foobarfoobarfoobarfoobarfoobafoobarfoobarfoobar"}
	jwtAuthenticator, err := utils.NewJWTAuthenticator(cfg.JWTSecret)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		body         models.MakeWithdrawalReq
		expectedCode int
		buildStubs   func(s *mocks.MockRepository)
		userID       int
	}{
		{
			name:         "Happy case",
			body:         models.MakeWithdrawalReq{Order: "4111111111111111", Sum: 99.0},
			expectedCode: http.StatusOK,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(1).
					Return(newFakeBalance(1, 100.0, 0), nil)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), "4111111111111111", 1, 99.0).
					Times(1).
					Return(nil)
			},
		},
		{
			name:         "Isufficient funds",
			body:         models.MakeWithdrawalReq{Order: "4111111111111111", Sum: 99.0},
			expectedCode: http.StatusPaymentRequired,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(1).
					Return(newFakeBalance(1, 95.50, 0), nil)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Invalid order number",
			body:         models.MakeWithdrawalReq{Order: "12345678902", Sum: 99.0},
			expectedCode: http.StatusUnprocessableEntity,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(0)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Empty order number",
			body:         models.MakeWithdrawalReq{Order: "", Sum: 99.0},
			expectedCode: http.StatusBadRequest,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(0)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Invalid sum",
			body:         models.MakeWithdrawalReq{Order: "4111111111111111", Sum: 0},
			expectedCode: http.StatusBadRequest,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(0)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Withdrawal already registered for order",
			body:         models.MakeWithdrawalReq{Order: "4111111111111111", Sum: 99.0},
			expectedCode: http.StatusConflict,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(1).
					Return(newFakeBalance(1, 100.0, 0), nil)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), "4111111111111111", 1, 99.0).
					Times(1).
					Return(errors.New("someError"))
				r.EXPECT().
					IsUniqueViolation(gomock.Any(), gomock.Any()).
					Times(1).
					Return(true)
			},
		},
		{
			name:         "Error on make withdrawal",
			body:         models.MakeWithdrawalReq{Order: "4111111111111111", Sum: 99.0},
			expectedCode: http.StatusInternalServerError,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetBalance(gomock.Any(), 1).
					Times(1).
					Return(newFakeBalance(1, 100.0, 0), nil)
				r.EXPECT().
					WithdrawFromBalance(gomock.Any(), "4111111111111111", 1, 99.0).
					Times(1).
					Return(errors.New("someError"))
				r.EXPECT().
					IsUniqueViolation(gomock.Any(), gomock.Any()).
					Times(1).
					Return(false)
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

			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(json))
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			s := service.NewGophermartService(repo, &cfg)
			handler := NewHandler(*s, jwtAuthenticator)

			handler.MakeWithdrawal(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}

		})
	}

}
