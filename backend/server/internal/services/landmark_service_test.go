// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLandmarkRepo struct{ mock.Mock }
func (m *mockLandmarkRepo) UpdateAllOwners(ctx context.Context, ownerClanID *models.Clan) error { return nil }
func (m *mockLandmarkRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Landmark, error) {
	args := m.Called(ctx, id)
	var l *models.Landmark
	if args.Get(0) != nil {
		l = args.Get(0).(*models.Landmark)
	}
	return l, args.Error(1)
}
func (m *mockLandmarkRepo) CheckDistance(ctx context.Context, landmarkID uuid.UUID, lat, lng float64) (bool, error) {
	args := m.Called(ctx, landmarkID, lat, lng)
	return args.Bool(0), args.Error(1)
}

type mockUserViolationRepo struct{ mock.Mock }
func (m *mockUserViolationRepo) Create(ctx context.Context, violation *models.UserViolation) error {
	return m.Called(ctx, violation).Error(0)
}
func (m *mockUserViolationRepo) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) { return 0, nil }
func (m *mockUserViolationRepo) BanUser(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockUserViolationRepo) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) { return nil, 0, nil }
func (m *mockUserViolationRepo) DeleteViolation(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockUserViolationRepo) HasActiveFakeGPSBan(ctx context.Context, userID uuid.UUID) (bool, *time.Time, error) {
	args := m.Called(ctx, userID)
	var t *time.Time
	if args.Get(1) != nil {
		t = args.Get(1).(*time.Time)
	}
	return args.Bool(0), t, args.Error(2)
}

type mockClanRepo struct{ mock.Mock }
func (m *mockClanRepo) Create(ctx context.Context, clan *models.Clan) error { return nil }
func (m *mockClanRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Clan, error) { return nil, nil }
func (m *mockClanRepo) AddScore(ctx context.Context, clanID uuid.UUID, score int) error {
	return m.Called(ctx, clanID, score).Error(0)
}
func (m *mockClanRepo) GetTopClan(ctx context.Context) (*models.Clan, error) { return nil, nil }
func (m *mockClanRepo) ResetWeeklyScores(ctx context.Context) error { return nil }
func (m *mockClanRepo) Update(ctx context.Context, clan *models.Clan) error { return nil }
func (m *mockClanRepo) GetLeaderboard(ctx context.Context, limit int) ([]models.Clan, error) { return nil, nil }
func (m *mockClanRepo) HardDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockClanRepo) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockClanRepo) FindByName(ctx context.Context, name string) (*models.Clan, error) { return nil, nil }

type mockUserRepoForLandmark struct{ mock.Mock }
func (m *mockUserRepoForLandmark) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	var u *models.User
	if args.Get(0) != nil {
		u = args.Get(0).(*models.User)
	}
	return u, args.Error(1)
}
func (m *mockUserRepoForLandmark) Create(ctx context.Context, user *models.User) error { return nil }
func (m *mockUserRepoForLandmark) FindByPhone(ctx context.Context, phone string) (*models.User, error) { return nil, nil }
func (m *mockUserRepoForLandmark) Update(ctx context.Context, user *models.User) error { return nil }

func TestLandmarkService_CheckIn_Mocked(t *testing.T) {
	landmarkRepo := new(mockLandmarkRepo)
	violationRepo := new(mockUserViolationRepo)
	clanRepo := new(mockClanRepo)
	userRepo := new(mockUserRepoForLandmark)

	service := services.NewLandmarkService(landmarkRepo, violationRepo, clanRepo, userRepo)
	ctx := context.Background()
	userID := uuid.New()
	landmarkID := uuid.New()

	violationRepo.On("HasActiveFakeGPSBan", ctx, userID).Return(false, nil, nil)
	violationRepo.On("Create", ctx, mock.AnythingOfType("*models.UserViolation")).Return(nil)

	_, err := service.CheckIn(ctx, userID, landmarkID, 10.0, 106.0, true)
	
	assert.Error(t, err)
	assert.Equal(t, services.ErrFakeGPS, err)
	violationRepo.AssertExpectations(t)
}

func TestLandmarkService_CheckIn_AlreadyBanned(t *testing.T) {
	landmarkRepo := new(mockLandmarkRepo)
	violationRepo := new(mockUserViolationRepo)
	clanRepo := new(mockClanRepo)
	userRepo := new(mockUserRepoForLandmark)

	service := services.NewLandmarkService(landmarkRepo, violationRepo, clanRepo, userRepo)
	ctx := context.Background()
	userID := uuid.New()
	landmarkID := uuid.New()

	banUntil := time.Now().Add(24*time.Hour)
	violationRepo.On("HasActiveFakeGPSBan", ctx, userID).Return(true, &banUntil, nil)

	_, err := service.CheckIn(ctx, userID, landmarkID, 10.0, 106.0, false)
	
	assert.Error(t, err)
	assert.Equal(t, services.ErrFakeGPS, err)
}

func TestLandmarkService_CheckIn_OutOfRange(t *testing.T) {
	landmarkRepo := new(mockLandmarkRepo)
	violationRepo := new(mockUserViolationRepo)
	clanRepo := new(mockClanRepo)
	userRepo := new(mockUserRepoForLandmark)

	service := services.NewLandmarkService(landmarkRepo, violationRepo, clanRepo, userRepo)
	ctx := context.Background()
	userID := uuid.New()
	landmarkID := uuid.New()

	violationRepo.On("HasActiveFakeGPSBan", ctx, userID).Return(false, nil, nil)
	landmarkRepo.On("CheckDistance", ctx, landmarkID, 10.0, 106.0).Return(false, nil)

	_, err := service.CheckIn(ctx, userID, landmarkID, 10.0, 106.0, false)
	
	assert.Error(t, err)
	assert.Equal(t, services.ErrOutOfRange, err)
}

func TestLandmarkService_CheckIn_Success(t *testing.T) {
	landmarkRepo := new(mockLandmarkRepo)
	violationRepo := new(mockUserViolationRepo)
	clanRepo := new(mockClanRepo)
	userRepo := new(mockUserRepoForLandmark)

	service := services.NewLandmarkService(landmarkRepo, violationRepo, clanRepo, userRepo)
	ctx := context.Background()
	userID := uuid.New()
	landmarkID := uuid.New()
	clanID := uuid.New()

	violationRepo.On("HasActiveFakeGPSBan", ctx, userID).Return(false, nil, nil)
	landmarkRepo.On("CheckDistance", ctx, landmarkID, 10.0, 106.0).Return(true, nil)
	landmark := &models.Landmark{ID: landmarkID}
	landmarkRepo.On("FindByID", ctx, landmarkID).Return(landmark, nil)
	
	user := &models.User{ID: userID, ClanID: &clanID}
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	clanRepo.On("AddScore", ctx, clanID, 10).Return(nil)

	res, err := service.CheckIn(ctx, userID, landmarkID, 10.0, 106.0, false)
	
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, landmarkID, res.ID)
	clanRepo.AssertExpectations(t)
}
