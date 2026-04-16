package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type adminUsageConversationRepoStub struct {
	service.UsageLogRepository

	conversation *service.UsageLogConversation
	err          error
	capturedID   int64
}

func (s *adminUsageConversationRepoStub) ReplaceRequestConversation(ctx context.Context, requestID string, userID, apiKeyID int64, snapshots []service.RequestConversationSnapshot) error {
	return nil
}

func (s *adminUsageConversationRepoStub) GetRequestConversationByUsageLogID(ctx context.Context, usageLogID int64) (*service.UsageLogConversation, error) {
	s.capturedID = usageLogID
	if s.err != nil {
		return nil, s.err
	}
	return s.conversation, nil
}

func newAdminUsageConversationTestRouter(repo *adminUsageConversationRepoStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	usageSvc := service.NewUsageService(repo, nil, nil, nil)
	handler := NewUsageHandler(usageSvc, nil, nil, nil)
	router := gin.New()
	router.GET("/admin/usage/:id/request-conversation", handler.GetRequestConversation)
	return router
}

func TestAdminUsageGetRequestConversationSuccess(t *testing.T) {
	repo := &adminUsageConversationRepoStub{
		conversation: &service.UsageLogConversation{
			UsageLogID: 123,
			RequestID:  "req_123",
			UserID:     7,
			APIKeyID:   11,
			Snapshots: []service.RequestConversationSnapshot{
				{Sequence: 1, Stage: service.RequestConversationStageInbound, Payload: `{"messages":[{"role":"user","content":"hello"}]}`},
				{Sequence: 2, Stage: service.RequestConversationStageUpstreamFinal, Payload: `{"messages":[{"role":"user","content":"hello"}],"model":"claude-3-7-sonnet"}`},
			},
		},
	}
	router := newAdminUsageConversationTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/admin/usage/123/request-conversation", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(123), repo.capturedID)

	var payload struct {
		Code int                          `json:"code"`
		Data service.UsageLogConversation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.Equal(t, 0, payload.Code)
	require.Equal(t, int64(123), payload.Data.UsageLogID)
	require.Len(t, payload.Data.Snapshots, 2)
	require.Equal(t, service.RequestConversationStageInbound, payload.Data.Snapshots[0].Stage)
	require.Equal(t, service.RequestConversationStageUpstreamFinal, payload.Data.Snapshots[1].Stage)
}

func TestAdminUsageGetRequestConversationInvalidID(t *testing.T) {
	repo := &adminUsageConversationRepoStub{}
	router := newAdminUsageConversationTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/admin/usage/not-a-number/request-conversation", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminUsageGetRequestConversationNotFound(t *testing.T) {
	repo := &adminUsageConversationRepoStub{
		err: service.ErrUsageLogConversationNotFound,
	}
	router := newAdminUsageConversationTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/admin/usage/99/request-conversation", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
