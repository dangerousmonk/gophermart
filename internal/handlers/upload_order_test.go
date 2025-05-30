package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/service/mocks"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newFakeOrder(userID int, number string, status models.OrderStatus) models.Order {
	return models.Order{
		ID:         1,
		Number:     number,
		UserID:     userID,
		Status:     status,
		Accrual:    0,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Active:     true,
	}
}

func TestUploadOrderHandler(t *testing.T) {
	cfg := config.Config{ServerAddr: "http://localhost:8080", JWTSecret: "foobarfoobarfoobarfoobarfoobafoobarfoobarfoobar"}
	jwtAuthenticator, err := utils.NewJWTAuthenticator(cfg.JWTSecret)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		orderNum     string
		expectedCode int
		buildStubs   func(s *mocks.MockRepository)
		userID       int
	}{
		{
			name:         "Happy case",
			orderNum:     `4111111111111111`,
			expectedCode: http.StatusAccepted,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), "4111111111111111").
					Times(1).
					Return(models.Order{}, sql.ErrNoRows)

				r.EXPECT().
					UploadOrder(gomock.Any(), "4111111111111111", 1, models.StatusNew).
					Times(1).
					Return(int64(1), nil)
			},
		},
		{
			name:         "Bad order number",
			orderNum:     ``,
			expectedCode: http.StatusBadRequest,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), gomock.Any()).
					Times(0)

				r.EXPECT().
					UploadOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Not valid order number",
			orderNum:     `12345678902`,
			expectedCode: http.StatusUnprocessableEntity,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), gomock.Any()).
					Times(0)

				r.EXPECT().
					UploadOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Order already uploaded by user",
			orderNum:     `4111111111111111`,
			expectedCode: http.StatusOK,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), "4111111111111111").
					Times(1).
					Return(newFakeOrder(1, "4111111111111111", models.StatusNew), nil)

				r.EXPECT().
					UploadOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Order already uploaded by different user",
			orderNum:     `4111111111111111`,
			expectedCode: http.StatusConflict,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), "4111111111111111").
					Times(1).
					Return(newFakeOrder(2, "4111111111111111", models.StatusNew), nil)

				r.EXPECT().
					UploadOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:         "Some upload error",
			orderNum:     `4111111111111111`,
			expectedCode: http.StatusInternalServerError,
			userID:       1,
			buildStubs: func(r *mocks.MockRepository) {
				r.EXPECT().
					GetOrderByNumber(gomock.Any(), "4111111111111111").
					Times(1).
					Return(models.Order{}, sql.ErrNoRows)

				r.EXPECT().
					UploadOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(int64(0), errors.New("someError"))
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

			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(tc.orderNum))
			ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, tc.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			s := service.NewGophermartService(repo, &cfg)
			handler := NewHandler(*s, jwtAuthenticator)

			handler.UploadOrder(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			}

		})
	}

}
