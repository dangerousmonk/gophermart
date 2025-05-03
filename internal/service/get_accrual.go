package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

const sleepTime = time.Second * 2

func (s *GophermartService) GetAccrual(orderNumber string) (*models.AccrualExternal, error) {
	var accrual models.AccrualExternal
	url := fmt.Sprintf("%s%s%s", s.Cfg.AccrualAddr, "/api/user/orders/", orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("GetAccrual request to accrual service failed ", slog.Any("err", err))
		return nil, fmt.Errorf("get request to accrual failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &accrual, nil
	case http.StatusNoContent:
		slog.Info("GetAccrual order not registered in system ", slog.String("orderNumber", orderNumber))
		return nil, nil
	case http.StatusTooManyRequests:
		timeSleep, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			time.Sleep(sleepTime)
		} else {
			time.Sleep(time.Duration(timeSleep) * time.Second)
		}
		return s.GetAccrual(orderNumber)
	default:
		slog.Error("GetAccrual server returned unexpected status", slog.Int("statusCode", resp.StatusCode), slog.String("URL", url))
		return nil, fmt.Errorf("accrual server error: %s", resp.Status)
	}
}
