package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestIsUnauthorizedError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{name: "http status prefix", msg: "HTTP 401: upstream rejected key", want: true},
		{name: "upstream 401", msg: "upstream 401", want: true},
		{name: "unauthorized", msg: "Unauthorized", want: true},
		{name: "invalid api key underscore", msg: "invalid_api_key", want: true},
		{name: "invalid api key spaces", msg: "invalid api key", want: true},
		{name: "authentication failed", msg: "authentication failed", want: true},
		{name: "401 at start", msg: "401 from upstream", want: true},
		{name: "401 in middle", msg: "upstream returned 401 unauthorized", want: true},
		{name: "401 at end", msg: "upstream returned 401", want: true},
		{name: "4012 not status", msg: "got code 4012 unrelated", want: false},
		{name: "network timeout", msg: "network timeout", want: false},
		{name: "rate limited", msg: "rate limited 429", want: false},
		{name: "empty", msg: "", want: false},
		{name: "whitespace", msg: " \t\n ", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isUnauthorizedError(tt.msg))
		})
	}
}

func TestScheduledTestRunnerRunOnePlanAutoDisableOnUnauth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		planAutoDisable    bool
		resultStatus       string
		resultErrorMessage string
		wantDisableCalls   int
	}{
		{
			name:               "disables schedulable account on unauthorized failure",
			planAutoDisable:    true,
			resultStatus:       "failed",
			resultErrorMessage: "HTTP 401",
			wantDisableCalls:   1,
		},
		{
			name:               "does not disable when toggle is off",
			planAutoDisable:    false,
			resultStatus:       "failed",
			resultErrorMessage: "HTTP 401",
			wantDisableCalls:   0,
		},
		{
			name:               "does not disable successful result",
			planAutoDisable:    true,
			resultStatus:       "success",
			resultErrorMessage: "HTTP 401",
			wantDisableCalls:   0,
		},
		{
			name:               "does not disable non unauthorized failure",
			planAutoDisable:    true,
			resultStatus:       "failed",
			resultErrorMessage: "rate limited 429",
			wantDisableCalls:   0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accountID := int64(42)
			planRepo := &scheduledTestPlanRepoStub{}
			resultRepo := &scheduledTestResultRepoStub{}
			accountRepo := &scheduledRunnerAccountRepoStub{}
			tester := &scheduledAccountTesterStub{
				result: &ScheduledTestResult{
					Status:       tt.resultStatus,
					ErrorMessage: tt.resultErrorMessage,
					StartedAt:    time.Now(),
					FinishedAt:   time.Now(),
				},
			}

			runner := &ScheduledTestRunnerService{
				planRepo:       planRepo,
				scheduledSvc:   NewScheduledTestService(planRepo, resultRepo),
				accountTestSvc: tester,
				accountRepo:    accountRepo,
			}

			runner.runOnePlan(context.Background(), &ScheduledTestPlan{
				ID:                  7,
				AccountID:           accountID,
				ModelID:             "test-model",
				CronExpression:      "*/5 * * * *",
				MaxResults:          10,
				AutoDisableOnUnauth: tt.planAutoDisable,
			})

			require.Equal(t, 1, tester.calls)
			require.Equal(t, 1, resultRepo.createCalls)
			require.Equal(t, 1, resultRepo.pruneCalls)
			require.Equal(t, 1, planRepo.updateAfterRunCalls)
			require.Equal(t, tt.wantDisableCalls, accountRepo.setSchedulableCalls)
			if tt.wantDisableCalls > 0 {
				require.Equal(t, accountID, accountRepo.lastSetSchedulableID)
				require.False(t, accountRepo.lastSchedulable)
			}
		})
	}
}

type scheduledAccountTesterStub struct {
	result *ScheduledTestResult
	err    error
	calls  int
}

func (s *scheduledAccountTesterStub) RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error) {
	s.calls++
	return s.result, s.err
}

type scheduledTestPlanRepoStub struct {
	updateAfterRunCalls int
	createdPlans        []*ScheduledTestPlan
}

func (s *scheduledTestPlanRepoStub) Create(ctx context.Context, plan *ScheduledTestPlan) (*ScheduledTestPlan, error) {
	s.createdPlans = append(s.createdPlans, plan)
	return plan, nil
}

func (s *scheduledTestPlanRepoStub) GetByID(ctx context.Context, id int64) (*ScheduledTestPlan, error) {
	return nil, nil
}

func (s *scheduledTestPlanRepoStub) ListByAccountID(ctx context.Context, accountID int64) ([]*ScheduledTestPlan, error) {
	return nil, nil
}

func (s *scheduledTestPlanRepoStub) ListDue(ctx context.Context, now time.Time) ([]*ScheduledTestPlan, error) {
	return nil, nil
}

func (s *scheduledTestPlanRepoStub) Update(ctx context.Context, plan *ScheduledTestPlan) (*ScheduledTestPlan, error) {
	return plan, nil
}

func (s *scheduledTestPlanRepoStub) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledTestPlanRepoStub) UpdateAfterRun(ctx context.Context, id int64, lastRunAt time.Time, nextRunAt time.Time) error {
	s.updateAfterRunCalls++
	return nil
}

type scheduledTestResultRepoStub struct {
	createCalls int
	pruneCalls  int
}

func (s *scheduledTestResultRepoStub) Create(ctx context.Context, result *ScheduledTestResult) (*ScheduledTestResult, error) {
	s.createCalls++
	return result, nil
}

func (s *scheduledTestResultRepoStub) ListByPlanID(ctx context.Context, planID int64, limit int) ([]*ScheduledTestResult, error) {
	return nil, nil
}

func (s *scheduledTestResultRepoStub) PruneOldResults(ctx context.Context, planID int64, keepCount int) error {
	s.pruneCalls++
	return nil
}

type scheduledRunnerAccountRepoStub struct {
	setSchedulableCalls  int
	lastSetSchedulableID int64
	lastSchedulable      bool
}

func (s *scheduledRunnerAccountRepoStub) Create(ctx context.Context, account *Account) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) GetByID(ctx context.Context, id int64) (*Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ExistsByID(ctx context.Context, id int64) (bool, error) {
	return true, nil
}

func (s *scheduledRunnerAccountRepoStub) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) FindByExtraField(ctx context.Context, key string, value any) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListCRSAccountIDs(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) Update(ctx context.Context, account *Account) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListAllWithFilters(ctx context.Context, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListOAuthRefreshCandidates(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) RevertProxyFallback(ctx context.Context, accountID int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) UpdateSessionWindowEnd(ctx context.Context, id int64, end time.Time) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) UpdateLastUsed(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ClearError(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	s.setSchedulableCalls++
	s.lastSetSchedulableID = id
	s.lastSchedulable = schedulable
	return nil
}

func (s *scheduledRunnerAccountRepoStub) AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error) {
	return 0, nil
}

func (s *scheduledRunnerAccountRepoStub) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) ListSchedulableUngroupedByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *scheduledRunnerAccountRepoStub) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) SetModelRateLimit(ctx context.Context, id int64, scope string, resetAt time.Time, models ...string) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ClearTempUnschedulable(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ClearRateLimit(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ClearModelRateLimits(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return 0, nil
}

func (s *scheduledRunnerAccountRepoStub) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ResetQuotaUsed(ctx context.Context, id int64) error {
	return nil
}

func (s *scheduledRunnerAccountRepoStub) ListShadowsByParent(ctx context.Context, parentID int64) ([]*Account, error) {
	return nil, nil
}

var _ AccountRepository = (*scheduledRunnerAccountRepoStub)(nil)
