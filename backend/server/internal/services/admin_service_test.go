// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserViolationRepo struct {
	BanUserFn         func(ctx context.Context, userID uuid.UUID) error
	GetViolationsFn   func(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error)
	DeleteViolationFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserViolationRepo) Create(ctx context.Context, violation *models.UserViolation) error {
	return nil
}
func (m *mockUserViolationRepo) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) {
	return 0, nil
}
func (m *mockUserViolationRepo) BanUser(ctx context.Context, userID uuid.UUID) error {
	if m.BanUserFn != nil {
		return m.BanUserFn(ctx, userID)
	}
	return nil
}
func (m *mockUserViolationRepo) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	if m.GetViolationsFn != nil {
		return m.GetViolationsFn(ctx, page, limit)
	}
	return nil, 0, nil
}
func (m *mockUserViolationRepo) DeleteViolation(ctx context.Context, id uuid.UUID) error {
	if m.DeleteViolationFn != nil {
		return m.DeleteViolationFn(ctx, id)
	}
	return nil
}

type mockCourtRepoAdmin struct {
	UpdateCaseStatusFn func(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error
	GetCaseByIDFn     func(ctx context.Context, id uuid.UUID) (*models.CourtCase, error)
}

func (m *mockCourtRepoAdmin) CreateCase(ctx context.Context, courtCase *models.CourtCase) error { return nil }
func (m *mockCourtRepoAdmin) GetActiveCases(ctx context.Context, jurorID uuid.UUID, limit int) ([]models.CourtCase, error) { return nil, nil }
func (m *mockCourtRepoAdmin) GetExpiredVotingCases(ctx context.Context) ([]models.CourtCase, error) { return nil, nil }
func (m *mockCourtRepoAdmin) CreateVote(ctx context.Context, vote *models.CourtVote) error { return nil }
func (m *mockCourtRepoAdmin) HasUserVoted(ctx context.Context, caseID uuid.UUID, jurorID uuid.UUID) (bool, error) { return false, nil }
func (m *mockCourtRepoAdmin) CountVotesByCase(ctx context.Context, caseID uuid.UUID) (int64, int64, error) { return 0, 0, nil }

func (m *mockCourtRepoAdmin) UpdateCaseStatus(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error {
	if m.UpdateCaseStatusFn != nil {
		return m.UpdateCaseStatusFn(ctx, id, status)
	}
	return nil
}

func (m *mockCourtRepoAdmin) GetCaseByID(ctx context.Context, id uuid.UUID) (*models.CourtCase, error) {
	if m.GetCaseByIDFn != nil {
		return m.GetCaseByIDFn(ctx, id)
	}
	return &models.CourtCase{}, nil
}

type mockWalletRepoAdmin struct {
	UpdateBalanceFn func(ctx context.Context, userID uuid.UUID, amount float64) error
}
func (m *mockWalletRepoAdmin) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	if m.UpdateBalanceFn != nil {
		return m.UpdateBalanceFn(ctx, userID, delta)
	}
	return nil
}
func (m *mockWalletRepoAdmin) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) { return nil, nil }
func (m *mockWalletRepoAdmin) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error { return nil }
func (m *mockWalletRepoAdmin) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) { return false, nil }
func (m *mockWalletRepoAdmin) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error { return nil }
func (m *mockWalletRepoAdmin) HoldBalance(ctx context.Context, userID uuid.UUID, amount float64, lock bool) error { return nil }

type mockTxManagerAdmin struct{}
func (m *mockTxManagerAdmin) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestAdminService_GetViolations(t *testing.T) {
	mockRepo := &mockUserViolationRepo{
		GetViolationsFn: func(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
			return []models.UserViolation{{ID: uuid.New()}}, 1, nil
		},
	}
	service := NewAdminService(mockRepo, nil, nil, nil, nil)
	
	// Test normal case
	violations, total, err := service.GetViolations(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, int64(1), total)

	// Test boundary case: page < 1
	violations2, total2, err2 := service.GetViolations(context.Background(), 0, 10)
	assert.NoError(t, err2)
	assert.Len(t, violations2, 1)
	assert.Equal(t, int64(1), total2)

	// Test boundary case: limit < 1 or limit > 100
	violations3, total3, err3 := service.GetViolations(context.Background(), 1, 150)
	assert.NoError(t, err3)
	assert.Len(t, violations3, 1)
	assert.Equal(t, int64(1), total3)
}

func TestAdminService_BanUser(t *testing.T) {
	mockRepo := &mockUserViolationRepo{}
	service := NewAdminService(mockRepo, nil, nil, nil, nil)
	err := service.BanUser(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestAdminService_DeleteViolationMedia(t *testing.T) {
	mockRepo := &mockUserViolationRepo{}
	service := NewAdminService(mockRepo, nil, nil, nil, nil)
	err := service.DeleteViolationMedia(context.Background(), uuid.New(), "")
	assert.NoError(t, err)
}

func TestAdminService_OverrideCourtCase(t *testing.T) {
	mockRepo := &mockCourtRepoAdmin{
		GetCaseByIDFn: func(ctx context.Context, id uuid.UUID) (*models.CourtCase, error) {
			return &models.CourtCase{
				PlaintiffID: uuid.New(),
				DefendantID: uuid.New(),
			}, nil
		},
	}
	mockWallet := &mockWalletRepoAdmin{}
	mockTx := &mockTxManagerAdmin{}
	service := NewAdminService(nil, mockRepo, nil, mockWallet, mockTx)
	err := service.OverrideCourtCase(context.Background(), uuid.New(), "guilty")
	assert.NoError(t, err)

	err2 := service.OverrideCourtCase(context.Background(), uuid.New(), "innocent")
	assert.NoError(t, err2)
}

type mockUserViolationRepoError struct {
	mockUserViolationRepo
}
func (m *mockUserViolationRepoError) BanUser(ctx context.Context, userID uuid.UUID) error {
	return errors.New("db error")
}
func (m *mockUserViolationRepoError) DeleteViolation(ctx context.Context, id uuid.UUID) error {
	return errors.New("db error")
}

func TestAdminService_Errors(t *testing.T) {
	mockUserRepo := &mockUserViolationRepoError{}
	service := NewAdminService(mockUserRepo, nil, nil, nil, nil)
	
	err := service.BanUser(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	err = service.DeleteViolationMedia(context.Background(), uuid.New(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

type mockCourtRepoAdminError struct {
	mockCourtRepoAdmin
}
func (m *mockCourtRepoAdminError) UpdateCaseStatus(ctx context.Context, id uuid.UUID, status models.CourtCaseStatus) error {
	return errors.New("db error")
}
func (m *mockCourtRepoAdminError) GetCaseByID(ctx context.Context, id uuid.UUID) (*models.CourtCase, error) {
	return nil, errors.New("db error")
}

func TestAdminService_OverrideCourtCase_Error(t *testing.T) {
	mockCourtRepo := &mockCourtRepoAdminError{}
	mockTx := &mockTxManagerAdmin{}
	service := NewAdminService(nil, mockCourtRepo, nil, nil, mockTx)
	
	err := service.OverrideCourtCase(context.Background(), uuid.New(), "dismissed")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
