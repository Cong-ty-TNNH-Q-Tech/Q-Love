// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"errors"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/google/uuid"
)

type LandmarkService interface {
	CheckIn(ctx context.Context, userID, landmarkID uuid.UUID, lat, lng float64, isMocked bool) (*models.Landmark, error)
}

type landmarkService struct {
	landmarkRepo   repository.LandmarkRepository
	violationRepo  repository.UserViolationRepository
	clanRepo       repository.ClanRepository
	userRepo       repository.UserRepository
}

// ErrFakeGPS is returned when a mock location is detected
var ErrFakeGPS = errors.New("ERR_FAKE_GPS_DETECTED")
// ErrOutOfRange is returned when user is too far from the landmark
var ErrOutOfRange = errors.New("ERR_OUT_OF_RANGE")

func NewLandmarkService(
	landmarkRepo repository.LandmarkRepository,
	violationRepo repository.UserViolationRepository,
	clanRepo repository.ClanRepository,
	userRepo repository.UserRepository,
) LandmarkService {
	return &landmarkService{
		landmarkRepo:  landmarkRepo,
		violationRepo: violationRepo,
		clanRepo:      clanRepo,
		userRepo:      userRepo,
	}
}

func (s *landmarkService) CheckIn(ctx context.Context, userID, landmarkID uuid.UUID, lat, lng float64, isMocked bool) (*models.Landmark, error) {
	// 1. Check if user already has an active fake GPS ban
	hasBan, _, err := s.violationRepo.HasActiveFakeGPSBan(ctx, userID)
	if err != nil {
		return nil, err
	}
	if hasBan {
		return nil, ErrFakeGPS
	}

	// 2. If current request uses a mock location, apply a 7-day ban
	if isMocked {
		violation := &models.UserViolation{
			UserID:   userID,
			Type:     "fake_gps",
			Reason:   "App detected Mock Location usage",
			IsActive: true,
		}
		if err := s.violationRepo.Create(ctx, violation); err != nil {
			return nil, err
		}
		return nil, ErrFakeGPS
	}

	// 3. Verify physical distance using PostGIS
	isWithin, err := s.landmarkRepo.CheckDistance(ctx, landmarkID, lat, lng)
	if err != nil {
		return nil, err
	}
	if !isWithin {
		return nil, ErrOutOfRange
	}

	// 4. Retrieve landmark details
	landmark, err := s.landmarkRepo.FindByID(ctx, landmarkID)
	if err != nil {
		return nil, err
	}

	// 5. Add 10 points to the user's Clan
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.ClanID != nil {
		// +10 points to Clan
		if err := s.clanRepo.AddScore(ctx, *user.ClanID, 10); err != nil {
			return nil, err
		}
	}

	return landmark, nil
}
