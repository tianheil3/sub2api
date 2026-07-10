package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountCreatesDefaultScheduledTestPlan(t *testing.T) {
	t.Parallel()

	accountRepo := &createAccountScheduledRepoStub{}
	planRepo := &scheduledTestPlanRepoStub{}
	resultRepo := &scheduledTestResultRepoStub{}
	svc := &adminServiceImpl{
		accountRepo:      accountRepo,
		scheduledTestSvc: NewScheduledTestService(planRepo, resultRepo),
	}

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Name:                 "auto-scheduled",
		Platform:             PlatformOpenAI,
		Type:                 AccountTypeAPIKey,
		Credentials:          map[string]any{"api_key": "sk-test"},
		Concurrency:          1,
		Priority:             1,
		SkipDefaultGroupBind: true,
	})

	require.NoError(t, err)
	require.Equal(t, int64(123), account.ID)
	require.Len(t, planRepo.createdPlans, 1)

	plan := planRepo.createdPlans[0]
	require.Equal(t, account.ID, plan.AccountID)
	require.Empty(t, plan.ModelID)
	require.Equal(t, defaultScheduledTestCronExpression, plan.CronExpression)
	require.True(t, plan.Enabled)
	require.Equal(t, defaultScheduledTestMaxResults, plan.MaxResults)
	require.True(t, plan.AutoRecover)
	require.True(t, plan.AutoDisableOnUnauth)
	require.NotNil(t, plan.NextRunAt)
}

type createAccountScheduledRepoStub struct {
}

func (s *createAccountScheduledRepoStub) Create(ctx context.Context, account *Account) error {
	account.ID = 123
	return nil
}

func (s *createAccountScheduledRepoStub) GetByID(ctx context.Context, id int64) (*Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ExistsByID(ctx context.Context, id int64) (bool, error) {
	return true, nil
}

func (s *createAccountScheduledRepoStub) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) FindByExtraField(ctx context.Context, key string, value any) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListCRSAccountIDs(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) Update(ctx context.Context, account *Account) error {
	return nil
}

func (s *createAccountScheduledRepoStub) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}

func (s *createAccountScheduledRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}

func (s *createAccountScheduledRepoStub) ListAllWithFilters(ctx context.Context, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListOAuthRefreshCandidates(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) UpdateSessionWindowEnd(ctx context.Context, id int64, end time.Time) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) UpdateLastUsed(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
}

func (s *createAccountScheduledRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ClearError(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
}

func (s *createAccountScheduledRepoStub) AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error) {
	return 0, nil
}

func (s *createAccountScheduledRepoStub) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) ListSchedulableUngroupedByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}

func (s *createAccountScheduledRepoStub) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
}

func (s *createAccountScheduledRepoStub) SetModelRateLimit(ctx context.Context, id int64, scope string, resetAt time.Time, models ...string) error {
	return nil
}

func (s *createAccountScheduledRepoStub) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
}

func (s *createAccountScheduledRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ClearTempUnschedulable(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ClearRateLimit(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ClearModelRateLimits(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
}

func (s *createAccountScheduledRepoStub) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
}

func (s *createAccountScheduledRepoStub) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return int64(len(ids)), nil
}

func (s *createAccountScheduledRepoStub) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ResetQuotaUsed(ctx context.Context, id int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) RevertProxyFallback(ctx context.Context, accountID int64) error {
	return nil
}

func (s *createAccountScheduledRepoStub) ListShadowsByParent(ctx context.Context, parentID int64) ([]*Account, error) {
	return nil, nil
}

var _ AccountRepository = (*createAccountScheduledRepoStub)(nil)
