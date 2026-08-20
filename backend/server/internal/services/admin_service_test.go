// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
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

type mockCourtCaseRepo struct {
	UpdateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
}

func (m *mockCourtCaseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status)
	}
	return nil
}

func TestAdminService_GetViolations(t *testing.T) {
	mockRepo := &mockUserViolationRepo{
		GetViolationsFn: func(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
			return []models.UserViolation{{ID: uuid.New()}}, 1, nil
		},
	}
	service := NewAdminService(mockRepo, nil, nil)
	
	// Test pagination bounds
	v, total, err := service.GetViolations(context.Background(), -1, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, v, 1)
}

func TestAdminService_BanUser(t *testing.T) {
	mockRepo := &mockUserViolationRepo{}
	service := NewAdminService(mockRepo, nil, nil)
	err := service.BanUser(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestAdminService_DeleteViolationMedia(t *testing.T) {
	mockRepo := &mockUserViolationRepo{}
	service := NewAdminService(mockRepo, nil, nil)
	err := service.DeleteViolationMedia(context.Background(), uuid.New(), "")
	assert.NoError(t, err)
}

func TestAdminService_OverrideCourtCase(t *testing.T) {
	mockRepo := &mockCourtCaseRepo{}
	service := NewAdminService(nil, mockRepo, nil)
	err := service.OverrideCourtCase(context.Background(), uuid.New(), "dismissed")
	assert.NoError(t, err)
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
	service := NewAdminService(mockUserRepo, nil, nil)
	
	err := service.BanUser(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	err = service.DeleteViolationMedia(context.Background(), uuid.New(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

type mockCourtCaseRepoError struct {
	mockCourtCaseRepo
}
func (m *mockCourtCaseRepoError) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return errors.New("db error")
}

func TestAdminService_OverrideCourtCase_Error(t *testing.T) {
	mockCourtRepo := &mockCourtCaseRepoError{}
	service := NewAdminService(nil, mockCourtRepo, nil)
	
	err := service.OverrideCourtCase(context.Background(), uuid.New(), "dismissed")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
