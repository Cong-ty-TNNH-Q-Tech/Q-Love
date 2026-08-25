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
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockUserViolationRepoCourt struct {
	mock.Mock
}

func (m *mockUserViolationRepoCourt) Create(ctx context.Context, violation *models.UserViolation) error {
	args := m.Called(ctx, violation)
	return args.Error(0)
}

func (m *mockUserViolationRepoCourt) DeleteViolation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserViolationRepoCourt) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]models.UserViolation), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserViolationRepoCourt) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) {
	return 0, nil
}

func (m *mockUserViolationRepoCourt) BanUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestCourtWorker_EvaluateExpiredCases_Guilty(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockViolation := new(mockUserViolationRepoCourt)
	logger := zap.NewNop()

	worker := NewCourtWorker(mockCourt, mockViolation, nil, logger, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	caseID := uuid.New()
	defendantID := uuid.New()

	cases := []models.CourtCase{
		{
			ID:          caseID,
			DefendantID: defendantID,
			Status:      models.CourtCaseStatusVoting,
			ExpiresAt:   time.Now().Add(-1 * time.Hour), // expired
		},
	}

	mockCourt.On("GetExpiredVotingCases", mock.Anything).Return(cases, nil)
	// 50 votes total, 35 guilty (70%)
	mockCourt.On("CountVotesByCase", mock.Anything, caseID).Return(int64(50), int64(35), nil)

	mockViolation.On("Create", mock.Anything, mock.AnythingOfType("*models.UserViolation")).Return(nil)
	mockViolation.On("BanUser", mock.Anything, defendantID).Return(nil)

	mockCourt.On("UpdateCaseStatus", mock.Anything, caseID, models.CourtCaseStatusGuilty).Return(nil)

	worker.evaluateExpiredCases(context.Background())

	mockCourt.AssertExpectations(t)
	mockViolation.AssertExpectations(t)
}

func TestCourtWorker_EvaluateExpiredCases_NotGuilty(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockViolation := new(mockUserViolationRepoCourt)
	logger := zap.NewNop()

	worker := NewCourtWorker(mockCourt, mockViolation, nil, logger, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	caseID := uuid.New()
	defendantID := uuid.New()

	cases := []models.CourtCase{
		{
			ID:          caseID,
			DefendantID: defendantID,
			Status:      models.CourtCaseStatusVoting,
			ExpiresAt:   time.Now().Add(-1 * time.Hour), // expired
		},
	}

	mockCourt.On("GetExpiredVotingCases", mock.Anything).Return(cases, nil)
	// 50 votes total, 30 guilty (60% -> not guilty, requires > 65%)
	mockCourt.On("CountVotesByCase", mock.Anything, caseID).Return(int64(50), int64(30), nil)

	// Violation repo should NOT be called
	
	mockCourt.On("UpdateCaseStatus", mock.Anything, caseID, models.CourtCaseStatusNotGuilty).Return(nil)

	worker.evaluateExpiredCases(context.Background())

	mockCourt.AssertExpectations(t)
	mockViolation.AssertExpectations(t)
}

func TestCourtWorker_Start(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockViolation := new(mockUserViolationRepoCourt)
	logger := zap.NewNop()

	worker := NewCourtWorker(mockCourt, mockViolation, nil, logger, &mockWalletRepoCourt{}, &mockTxManagerCourt{})
	ctx, cancel := context.WithCancel(context.Background())
	
	worker.Start(ctx)
	// Give it a moment to start the goroutines
	time.Sleep(100 * time.Millisecond)
	
	cancel()
	// Give it a moment to stop
	time.Sleep(100 * time.Millisecond)
}
func TestCourtWorker_EvaluateExpiredCases_RepoErrors(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockViolation := new(mockUserViolationRepoCourt)
	logger := zap.NewNop()

	worker := NewCourtWorker(mockCourt, mockViolation, nil, logger, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	// 1. GetExpiredVotingCases returns error
	mockCourt.On("GetExpiredVotingCases", mock.Anything).Return([]models.CourtCase{}, errors.New("db error")).Once()
	worker.evaluateExpiredCases(context.Background())

	// 2. CountVotesByCase returns error
	caseID := uuid.New()
	cases := []models.CourtCase{{ID: caseID, Status: models.CourtCaseStatusVoting}}
	mockCourt.On("GetExpiredVotingCases", mock.Anything).Return(cases, nil).Once()
	mockCourt.On("CountVotesByCase", mock.Anything, caseID).Return(int64(0), int64(0), errors.New("count error")).Once()
	worker.evaluateExpiredCases(context.Background())
}

func TestCourtWorker_EvaluateExpiredCases_ViolationErrors(t *testing.T) {
	mockCourt := new(mockCourtRepo)
	mockViolation := new(mockUserViolationRepoCourt)
	logger := zap.NewNop()

	worker := NewCourtWorker(mockCourt, mockViolation, nil, logger, &mockWalletRepoCourt{}, &mockTxManagerCourt{})

	caseID := uuid.New()
	cases := []models.CourtCase{{ID: caseID, DefendantID: uuid.New(), Status: models.CourtCaseStatusVoting}}
	
	// Create violation error
	mockCourt.On("GetExpiredVotingCases", mock.Anything).Return(cases, nil).Once()
	mockCourt.On("CountVotesByCase", mock.Anything, caseID).Return(int64(50), int64(35), nil).Once()
	mockViolation.On("Create", mock.Anything, mock.Anything).Return(errors.New("create err")).Once()
	mockViolation.On("BanUser", mock.Anything, mock.Anything).Return(errors.New("ban err")).Once()
	mockCourt.On("UpdateCaseStatus", mock.Anything, caseID, models.CourtCaseStatusGuilty).Return(errors.New("update err")).Once()
	
	worker.evaluateExpiredCases(context.Background())
}
