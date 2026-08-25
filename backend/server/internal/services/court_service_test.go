// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks

type mockCourtRepo struct {
	mock.Mock
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

func (m *mockCourtRepo) GetActiveCases(ctx context.Context, limit int) ([]models.CourtCase, error) {
	args := m.Called(ctx, limit)
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

// Tests

func TestCourtService_FileLawsuit_Success(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)

	service := NewCourtService(mockCourt, mockMatch, nil)

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

	service := NewCourtService(mockCourt, mockMatch, nil)

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

	service := NewCourtService(mockCourt, mockMatch, nil)

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
	svc := NewCourtService(mockCourt, mockMatch, nil)

	activeCases := []models.CourtCase{
		{ID: uuid.New(), Status: models.CourtCaseStatusVoting},
	}
	mockCourt.On("GetActiveCases", mock.Anything, 10).Return(activeCases, nil)

	cases, err := svc.GetFeed(context.Background(), uuid.New(), 10)
	assert.NoError(t, err)
	assert.Len(t, cases, 1)
}

func TestCourtService_WithdrawCase_Success(t *testing.T) {
	caseID := uuid.New()
	plaintiffID := uuid.New()
	
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil)

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
	svc := NewCourtService(mockCourt, mockMatch, nil)

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
	svc := NewCourtService(mockCourt, mockMatch, nil)

	userID := uuid.New()
	matchID := uuid.New()

	_, err := svc.FileLawsuit(context.Background(), userID, userID, matchID, "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "cannot sue yourself", err.Error())
}
func TestCourtService_FileLawsuit_MatchNotFound(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil)
	
	mockMatch.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	
	_, err := svc.FileLawsuit(context.Background(), uuid.New(), uuid.New(), uuid.New(), "Ghosting")
	assert.Error(t, err)
	assert.Equal(t, "match not found", err.Error())
}

func TestCourtService_FileLawsuit_Unauthorized(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockMatch := new(mockMatchRepoForCourt)
	svc := NewCourtService(mockCourt, mockMatch, nil)
	
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
	svc := NewCourtService(mockCourt, mockMatch, nil)
	
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
	svc := NewCourtService(mockCourt, nil, nil)
	
	courtCase := &models.CourtCase{
		Status: models.CourtCaseStatusSettled,
	}
	mockCourt.On("GetCaseByID", mock.Anything, mock.Anything).Return(courtCase, nil)
	
	err := svc.VoteCase(context.Background(), uuid.New(), uuid.New(), models.CourtVoteGuilty)
	assert.Error(t, err)
	assert.Equal(t, "voting is closed for this case", err.Error())
}
