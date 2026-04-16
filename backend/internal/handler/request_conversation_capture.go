package handler

import (
	"bytes"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

func buildRequestConversationSnapshots(c *gin.Context, originalBody, currentBody []byte, account *service.Account) []service.RequestConversationSnapshot {
	snapshots := make([]service.RequestConversationSnapshot, 0, 4)
	if len(originalBody) > 0 {
		snapshots = append(snapshots, service.RequestConversationSnapshot{
			Stage:   service.RequestConversationStageInbound,
			Kind:    "client_request",
			Payload: string(originalBody),
		})
	}
	if len(currentBody) > 0 && !bytes.Equal(originalBody, currentBody) {
		snapshot := service.RequestConversationSnapshot{
			Stage:   service.RequestConversationStageGatewayRewrite,
			Kind:    "gateway_mutation",
			Payload: string(currentBody),
		}
		if account != nil {
			snapshot.Platform = strings.TrimSpace(account.Platform)
			snapshot.AccountID = &account.ID
			snapshot.AccountName = strings.TrimSpace(account.Name)
		}
		snapshots = append(snapshots, snapshot)
	}
	if c != nil {
		if raw, ok := c.Get(service.OpsUpstreamErrorsKey); ok {
			if events, ok := raw.([]*service.OpsUpstreamErrorEvent); ok {
				for _, event := range events {
					if event == nil || strings.TrimSpace(event.UpstreamRequestBody) == "" {
						continue
					}
					snapshot := service.RequestConversationSnapshot{
						Stage:       service.RequestConversationStageUpstreamRetry,
						Kind:        strings.TrimSpace(event.Kind),
						Payload:     strings.TrimSpace(event.UpstreamRequestBody),
						Platform:    strings.TrimSpace(event.Platform),
						AccountName: strings.TrimSpace(event.AccountName),
						UpstreamURL: strings.TrimSpace(event.UpstreamURL),
					}
					if event.AccountID > 0 {
						accountID := event.AccountID
						snapshot.AccountID = &accountID
					}
					snapshots = append(snapshots, snapshot)
				}
			}
		}
		if raw, ok := c.Get(service.OpsUpstreamRequestBodyKey); ok {
			payload := ""
			switch body := raw.(type) {
			case string:
				payload = strings.TrimSpace(body)
			case []byte:
				payload = strings.TrimSpace(string(body))
			}
			if payload != "" {
				snapshot := service.RequestConversationSnapshot{
					Stage:   service.RequestConversationStageUpstreamFinal,
					Kind:    "success",
					Payload: payload,
				}
				if account != nil {
					snapshot.Platform = strings.TrimSpace(account.Platform)
					snapshot.AccountID = &account.ID
					snapshot.AccountName = strings.TrimSpace(account.Name)
				}
				snapshots = append(snapshots, snapshot)
			}
		}
	}
	for idx := range snapshots {
		snapshots[idx].Sequence = idx + 1
	}
	return snapshots
}
