// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"database/sql"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// Mocks

type mockCourtRepo struct {
	mock.Mock
}

type mockWalletRepoCourt struct{}
func (m *mockWalletRepoCourt) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error { return nil }
func (m *mockWalletRepoCourt) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) { return nil, nil }
func (m *mockWalletRepoCourt) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockWalletRepoCourt) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }
func (m *mockWalletRepoCourt) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error { return nil }
func (m *mockWalletRepoCourt) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockWalletRepoCourt) ReleaseHoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }

type mockTxManagerCourt struct{}
func (m *mockTxManagerCourt) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func (m *mockCourtRepo) CreateCase(ctx context.Context, courtCase *models.CourtCase) error {
	args := m.Called(ctx, courtCase)
	return args.Error(0)
}

func (m *mockCourtRepo) GetCaseByID(ctx context.Context, id uuid.UUID) (*models.CourtCase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CourtCase), args.Error(1)
}

func (m *mockCourtRepo) GetActiveCases(ctx context.Context, jurorID uuid.UUID, limit int) ([]models.CourtCase, error) {
	args := m.Called(ctx, jurorID, limit)
	return args.Get(0).([]models.CourtCase), args.Error(1)
}

func (m *mockCourtRepo) GetExpiredVotingCases(ctx context.Context) ([]models.CourtCase, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.CourtCase), args.Error(1)
}

func (m *mockCourtRepo) UpdateCaseStatus(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockCourtRepo) CreateVote(ctx context.Context, vote *models.CourtVote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *mockCourtRepo) HasUserVoted(ctx context.Context, caseID uuid.UUID, jurorID uuid.UUID) (bool, error) {
	args := m.Called(ctx, caseID, jurorID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCourtRepo) CountVotesByCase(ctx context.Context, caseID uuid.UUID) (int64, int64, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

type mockMatchRepoForCourt struct {
	mock.Mock
}

func (m *mockMatchRepoForCourt) Create(ctx context.Context, match *models.Match) error {
	return nil
}

func (m *mockMatchRepoForCourt) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *mockMatchRepoForCourt) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *mockMatchRepoForCourt) SoftDelete(ctx context.Context, matchID uuid.UUID) error {
	args := m.Called(ctx, matchID)
	return args.Error(0)
}

func (m *mockMatchRepoForCourt) GetByUserIDs(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) {
	return nil, nil
}

func (m *mockMatchRepoForCourt) GetActiveMatchesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Match, error) {
	return nil, nil
}

func (m *mockMatchRepoForCourt) UpdateStreak(ctx context.Context, id uuid.UUID, score int) error {
	return nil
}

func (m *mockMatchRepoForCourt) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error {
	return nil
}

func (m *mockMatchRepoForCourt) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockMatchRepoForCourt) FindByUsers(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) {
	return nil, nil
}

// Tests

func TestCourtService_FileLawsuit_Success(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)

	service := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	plaintiffID := uuid.New()
	defendantID := uuid.New()
	matchID := uuid.New()

	match := &models.Match{
		ID:                matchID,
		User1ID:           plaintiffID,
		User2ID:           defendantID,
		StreakScore:       10,
		LastInteractionAt: time.Now().Add(-50 * time.Hour),
	}

	mockMatch.On("FindByID", mock.Anything, matchID).Return(match, nil)
	mockCourt.On("CreateCase", mock.Anything, mock.AnythingOfType("*models.CourtCase")).Return(nil)

	courtCase, err := service.FileLawsuit(context.Background(), plaintiffID, defendantID, matchID, "Ghosting")

	assert.NoError(t, err)
	assert.NotNil(t, courtCase)
	assert.Equal(t, plaintiffID, courtCase.PlaintiffID)
	assert.Equal(t, defendantID, courtCase.DefendantID)
	assert.Equal(t, "Ghosting", courtCase.Reason)
	assert.Equal(t, models.CourtCaseStatusVoting, courtCase.Status)

	mockMatch.AssertExpectations(t)
	mockCourt.AssertExpectations(t)
}

func TestCourtService_FileLawsuit_FailsIfNot48Hours(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)

	service := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	plaintiffID := uuid.New()
	defendantID := uuid.New()
	matchID := uuid.New()

	match := &models.Match{
		ID:                matchID,
		User1ID:           plaintiffID,
		User2ID:           defendantID,
		StreakScore:       10,
		LastInteractionAt: time.Now().Add(-10 * time.Hour), // Less than 48h
	}

	mockMatch.On("FindByID", mock.Anything, matchID).Return(match, nil)

	_, err := service.FileLawsuit(context.Background(), plaintiffID, defendantID, matchID, "Ghosting")

	assert.Error(t, err)
	assert.Equal(t, "must wait at least 48 hours of silence before filing a lawsuit", err.Error())
}

func TestCourtService_VoteCase_Success(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)

	service := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	caseID := uuid.New()
	jurorID := uuid.New()

	courtCase := &models.CourtCase{
		ID:          caseID,
		PlaintiffID: uuid.New(),
		DefendantID: uuid.New(),
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(10 * time.Hour),
	}

	mockCourt.On("GetCaseByID", mock.Anything, caseID).Return(courtCase, nil)
	mockCourt.On("HasUserVoted", mock.Anything, caseID, jurorID).Return(false, nil)
	mockCourt.On("CreateVote", mock.Anything, mock.AnythingOfType("*models.CourtVote")).Return(nil)

	err := service.VoteCase(context.Background(), caseID, jurorID, models.CourtVoteGuilty)

	assert.NoError(t, err)
	mockCourt.AssertExpectations(t)
}

func TestCourtService_GetFeed_Success(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	jurorID := uuid.New()

	activeCases := []models.CourtCase{
		{ID: uuid.New(), Status: models.CourtCaseStatusVoting},
	}
	mockCourt.On("GetActiveCases", mock.Anything, jurorID, 10).Return(activeCases, nil)

	cases, err := svc.GetFeed(context.Background(), jurorID, 10)
	assert.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestCourtService_WithdrawCase_Success(t *testing.T) {
	caseID := uuid.New()
	plaintiffID := uuid.New()
	
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	courtCase := &models.CourtCase{
		ID:          caseID,
		PlaintiffID: plaintiffID,
		Status:      models.CourtCaseStatusVoting,
	}
	mockCourt.On("GetCaseByID", mock.Anything, caseID).Return(courtCase, nil)
	mockCourt.On("UpdateCaseStatus", mock.Anything, caseID, models.CourtCaseStatusWithdrawn).Return(nil)

	err := svc.WithdrawCase(context.Background(), caseID, plaintiffID)
	assert.NoError(t, err)
}

func TestCourtService_WithdrawCase_NotPlaintiff(t *testing.T) {
	caseID := uuid.New()
	plaintiffID := uuid.New()
	wrongID := uuid.New()
	
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	courtCase := &models.CourtCase{
		ID:          caseID,
		PlaintiffID: plaintiffID,
		Status:      models.CourtCaseStatusVoting,
	}
	mockCourt.On("GetCaseByID", mock.Anything, caseID).Return(courtCase, nil)

	err := svc.WithdrawCase(context.Background(), caseID, wrongID)
	assert.Error(t, err)
	assert.Equal(t, "unauthorized to withdraw this case", err.Error())
}

func TestCourtService_FileLawsuit_CannotSueYourself(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	userID := uuid.New()
	matchID := uuid.New()

	_, err := svc.FileLawsuit(context.Background(), userID, userID, matchID, "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "cannot sue yourself", err.Error())
}
func TestCourtService_FileLawsuit_MatchNotFound(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	
	_, err := svc.FileLawsuit(context.Background(), uuid.New(), uuid.New(), uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "match not found", err.Error())
}

func TestCourtService_FileLawsuit_Unauthorized(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	match := &models.Match{
		User1ID: uuid.New(),
		User2ID: uuid.New(),
	}
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(match, nil)
	
	_, err := svc.FileLawsuit(context.Background(), uuid.New(), uuid.New(), uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "unauthorized to file lawsuit for this match", err.Error())
}

func TestCourtService_FileLawsuit_InvalidDefendant(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	match := &models.Match{
		User1ID: pID,
		User2ID: uuid.New(),
	}
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(match, nil)
	
	_, err := svc.FileLawsuit(context.Background(), pID, uuid.New(), uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "invalid defendant for this match", err.Error())
}

func TestCourtService_VoteCase_VotingClosed(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	courtCase := &models.CourtCase{
		Status: models.CourtCaseStatusSettled,
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	
	err := svc.VoteCase(context.Background(), uuid.New(), uuid.New(), models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "voting is closed for this case", err.Error())
}
func TestCourtService_FileLawsuit_StreakTooLow(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	dID := uuid.New()
	match := &models.Match{
		User1ID: pID,
		User2ID: dID,
		StreakScore: 5, // <= 5
	}
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(match, nil)
	
	_, err := svc.FileLawsuit(context.Background(), pID, dID, uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "streak must be greater than 5 to file a lawsuit", err.Error())
}

func TestCourtService_VoteCase_CaseNotFound(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	
	err := svc.VoteCase(context.Background(), uuid.New(), uuid.New(), models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "case not found", err.Error())
}

func TestCourtService_VoteCase_VotingExpired(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	courtCase := &models.CourtCase{
		Status: models.CourtCaseStatusVoting,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // expired
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	
	err := svc.VoteCase(context.Background(), uuid.New(), uuid.New(), models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "voting period has expired", err.Error())
}

func TestCourtService_VoteCase_PlaintiffDefendantCannotVote(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	courtCase := &models.CourtCase{
		Status: models.CourtCaseStatusVoting,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		PlaintiffID: pID,
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	
	err := svc.VoteCase(context.Background(), uuid.New(), pID, models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "plaintiff or defendant cannot vote in their own case", err.Error())
}

func TestCourtService_VoteCase_AlreadyVoted(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	courtCase := &models.CourtCase{
		Status: models.CourtCaseStatusVoting,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	mockCourt.On("HasUserVoted", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	
	err := svc.VoteCase(context.Background(), uuid.New(), uuid.New(), models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "already voted", err.Error())
}

func TestCourtService_WithdrawCase_CaseNotFound(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	
	err := svc.WithdrawCase(context.Background(), uuid.New(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, "case not found", err.Error())
}

func TestCourtService_WithdrawCase_NotInVotingPhase(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	courtCase := &models.CourtCase{
		PlaintiffID: pID,
		Status: models.CourtCaseStatusSettled,
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	
	err := svc.WithdrawCase(context.Background(), uuid.New(), pID)
	assert.Error(t, err)
	assert.Equal(t, "can only withdraw cases that are currently in voting phase", err.Error())
}

func TestCourtService_VoteCase_CreateVoteError(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	caseID := uuid.New()
	jurorID := uuid.New()
	
	courtCase := &models.CourtCase{
		ID:          caseID,
		Status:      models.CourtCaseStatusVoting,
		ExpiresAt:   time.Now().Add(10 * time.Hour),
	}
	
	mockCourt.On("GetCaseByID", mock.Anything, caseID).Return(courtCase, nil)
	mockCourt.On("HasUserVoted", mock.Anything, caseID, jurorID).Return(false, nil)
	mockCourt.On("CreateVote", mock.Anything, mock.AnythingOfType("*models.CourtVote")).Return(errors.New("db error"))
	
	err := svc.VoteCase(context.Background(), caseID, jurorID, models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCourtService_WithdrawCase_UpdateError(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	caseID := uuid.New()
	pID := uuid.New()
	courtCase := &models.CourtCase{
		PlaintiffID: pID,
		Status: models.CourtCaseStatusVoting,
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	mockCourt.On("UpdateCaseStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	
	err := svc.WithdrawCase(context.Background(), caseID, pID)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCourtService_FileLawsuit_TxFail(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	dID := uuid.New()
	match := &models.Match{
		User1ID: pID,
		User2ID: dID,
		StreakScore: 10,
		LastInteractionAt: time.Now().Add(-50 * time.Hour),
	}
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(match, nil)
	mockCourt.On("CreateCase", mock.Anything, mock.Anything).Return(errors.New("db error"))
	
	_, err := svc.FileLawsuit(context.Background(), pID, dID, uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCourtService_GetFeed_Error(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	svc := NewCourtService(mockCourt, nil, nil, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	
	mockCourt.On("GetActiveCases", mock.Anything, mock.Anything, mock.Anything).Return([]models.CourtCase{}, errors.New("db error"))
	
	_, err := svc.GetFeed(context.Background(), uuid.New(), 10)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

type mockWalletRepoCourtFail struct{}
func (m *mockWalletRepoCourtFail) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error { return errors.New("wallet err") }
func (m *mockWalletRepoCourtFail) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) { return nil, nil }
func (m *mockWalletRepoCourtFail) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockWalletRepoCourtFail) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }
func (m *mockWalletRepoCourtFail) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error { return nil }
func (m *mockWalletRepoCourtFail) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockWalletRepoCourtFail) ReleaseHoldBalance(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }

func TestCourtService_FileLawsuit_WalletFail(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil, &mockWalletRepoCourtFail{}, &mockTxManagerCourt{})
	
	pID := uuid.New()
	dID := uuid.New()
	match := &models.Match{
		User1ID: pID,
		User2ID: dID,
		StreakScore: 10,
		LastInteractionAt: time.Now().Add(-50 * time.Hour),
	}
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(match, nil)
	
	_, err := svc.FileLawsuit(context.Background(), pID, dID, uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "wallet err", err.Error())
}

func TestCourtService_FileLawsuit_WithRedis(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	mockWallet := new(mockWalletRepoCourt)
	
	svc := NewCourtService(mockCourt, mockMatch, client, mockWallet, &mockTxManagerCourt{})
	
	pID := uuid.New()
	dID := uuid.New()
	matchID := uuid.New()
	
	match := &models.Match{
		User1ID: pID,
		User2ID: dID,
		StreakScore: 10,
		LastInteractionAt: time.Now().Add(-50 * time.Hour),
	}
	mockMatch.On("FindByID", mock.Anything, matchID).Return(match, nil)
	mockCourt.On("CreateCase", mock.Anything, mock.Anything).Return(nil)
	
	caseData, err := svc.FileLawsuit(context.Background(), pID, dID, matchID, "Ghosting")
	assert.NoError(t, err)
	assert.NotNil(t, caseData)
}

func TestCourtService_FileLawsuit_WithRedisError(t *testing.T) {
	mr, _ := miniredis.Run()
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	mockWallet := new(mockWalletRepoCourt)
	
	svc := NewCourtService(mockCourt, mockMatch, client, mockWallet, &mockTxManagerCourt{})
	
	pID := uuid.New()
	dID := uuid.New()
	matchID := uuid.New()
	
	match := &models.Match{
		User1ID: pID,
		User2ID: dID,
		StreakScore: 10,
		LastInteractionAt: time.Now().Add(-50 * time.Hour),
	}
	mockMatch.On("FindByID", mock.Anything, matchID).Return(match, nil)
	mockCourt.On("CreateCase", mock.Anything, mock.Anything).Return(nil)
	
	// Close miniredis so XAdd fails!
	mr.Close()
	
	caseData, err := svc.FileLawsuit(context.Background(), pID, dID, matchID, "Ghosting")
	assert.NoError(t, err) // Lawsuit still succeeds even if redis fails!
	assert.NotNil(t, caseData)
}
