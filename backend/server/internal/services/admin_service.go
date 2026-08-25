// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/google/uuid"
)

type AdminService interface {
	GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error)
	BanUser(ctx context.Context, userID uuid.UUID) error
	DeleteViolationMedia(ctx context.Context, violationID uuid.UUID, objectKey string) error
	OverrideCourtCase(ctx context.Context, caseID uuid.UUID, status string) error
}

type adminService struct {
	violationRepo repository.UserViolationRepository
	courtCaseRepo repository.CourtCaseRepository
	r2Client      *storage.R2Client
	walletRepo    repository.WalletRepository
	txManager     repository.TransactionManager
}

func NewAdminService(
	violationRepo repository.UserViolationRepository,
	courtCaseRepo repository.CourtCaseRepository,
	r2Client *storage.R2Client,
	walletRepo repository.WalletRepository,
	txManager repository.TransactionManager,
) AdminService {
	return &adminService{
		violationRepo: violationRepo,
		courtCaseRepo: courtCaseRepo,
		r2Client:      r2Client,
		walletRepo:    walletRepo,
		txManager:     txManager,
	}
}

func (s *adminService) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.violationRepo.GetViolations(ctx, page, limit)
}

func (s *adminService) BanUser(ctx context.Context, userID uuid.UUID) error {
	return s.violationRepo.BanUser(ctx, userID)
}

func (s *adminService) DeleteViolationMedia(ctx context.Context, violationID uuid.UUID, objectKey string) error {
	// First delete from R2
	if s.r2Client != nil && objectKey != "" {
		if err := s.r2Client.DeleteObject(ctx, objectKey); err != nil {
			return err
		}
	}
	// Then mark violation as inactive
	return s.violationRepo.DeleteViolation(ctx, violationID)
}

func (s *adminService) OverrideCourtCase(ctx context.Context, caseID uuid.UUID, status string) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		courtCase, err := s.courtCaseRepo.FindByID(txCtx, caseID)
		if err != nil {
			return err
		}

		if err := s.courtCaseRepo.UpdateStatus(txCtx, caseID, status); err != nil {
			return err
		}

		courtFee := float64(50)
		if status == "guilty" {
			if err := s.walletRepo.UpdateBalance(txCtx, courtCase.DefendantID, -courtFee); err != nil {
				return err
			}
			if err := s.walletRepo.UpdateBalance(txCtx, courtCase.PlaintiffID, courtFee); err != nil {
				return err
			}
		} else if status == "innocent" {
			if err := s.walletRepo.UpdateBalance(txCtx, courtCase.PlaintiffID, -courtFee); err != nil {
				return err
			}
			if err := s.walletRepo.UpdateBalance(txCtx, courtCase.DefendantID, courtFee); err != nil {
				return err
			}
		}

		return nil
	})
}
