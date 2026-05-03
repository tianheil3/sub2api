package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type subscriptionResetRepoStub struct {
	service.UserSubscriptionRepository

	sub           *service.UserSubscription
	consumeCalled bool
}

func (r *subscriptionResetRepoStub) GetByID(_ context.Context, id int64) (*service.UserSubscription, error) {
	if r.sub == nil || r.sub.ID != id {
		return nil, service.ErrSubscriptionNotFound
	}
	cp := *r.sub
	return &cp, nil
}

func (r *subscriptionResetRepoStub) ConsumeOneDayAndResetDailyUsage(_ context.Context, id int64, newWindowStart time.Time) error {
	if r.sub == nil || r.sub.ID != id {
		return service.ErrSubscriptionNotFound
	}
	r.consumeCalled = true
	if !r.sub.ExpiresAt.AddDate(0, 0, -1).After(time.Now()) {
		return service.ErrAdjustWouldExpire
	}
	r.sub.DailyUsageUSD = 0
	r.sub.DailyWindowStart = &newWindowStart
	r.sub.ExpiresAt = r.sub.ExpiresAt.AddDate(0, 0, -1)
	return nil
}

func TestSubscriptionHandlerResetQuota(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dailyLimit := 100.0
	dailyWindowStart := time.Now().Add(-12 * time.Hour)
	originalExpiresAt := time.Now().Add(48 * time.Hour)

	repo := &subscriptionResetRepoStub{
		sub: &service.UserSubscription{
			ID:               501,
			UserID:           1,
			GroupID:          10,
			StartsAt:         time.Now().Add(-24 * time.Hour),
			ExpiresAt:        originalExpiresAt,
			Status:           service.SubscriptionStatusActive,
			DailyWindowStart: &dailyWindowStart,
			DailyUsageUSD:    12.34,
			Group: &service.Group{
				ID:               10,
				Name:             "Group One",
				SubscriptionType: service.SubscriptionTypeSubscription,
				Status:           service.StatusActive,
				DailyLimitUSD:    &dailyLimit,
			},
		},
	}

	svc := service.NewSubscriptionService(nil, repo, nil, nil, nil)
	h := NewSubscriptionHandler(svc)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1, Concurrency: 5})
		c.Next()
	})
	router.POST("/subscriptions/:id/reset-quota", h.ResetQuota)

	req := httptest.NewRequest(http.MethodPost, "/subscriptions/501/reset-quota", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, repo.consumeCalled)

	var payload struct {
		Code int `json:"code"`
		Data struct {
			DailyUsageUSD    float64    `json:"daily_usage_usd"`
			ExpiresAt        time.Time  `json:"expires_at"`
			DailyWindowStart *time.Time `json:"daily_window_start"`
			GroupID          int64      `json:"group_id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.Equal(t, 0, payload.Code)
	require.Equal(t, float64(0), payload.Data.DailyUsageUSD)
	require.NotNil(t, payload.Data.DailyWindowStart)
	require.Equal(t, 0, payload.Data.DailyWindowStart.Hour())
	require.True(t, payload.Data.ExpiresAt.Equal(originalExpiresAt.AddDate(0, 0, -1)))
	require.Equal(t, int64(10), payload.Data.GroupID)
}
