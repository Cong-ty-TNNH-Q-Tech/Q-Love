package services

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/google/uuid"
)

type mockMatchRepo struct {
	match *models.Match
	err   error
}

func (m *mockMatchRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	return m.match, m.err
}

func (m *mockMatchRepo) UpdateLastInteraction(ctx context.Context, id uuid.UUID, t time.Time) error {
	return nil
}

func (m *mockMatchRepo) Create(ctx context.Context, match *models.Match) error {
	return m.err
}

type mockChatRepo struct {
	err error
}

func (m *mockChatRepo) Create(ctx context.Context, message *models.ChatMessage) error {
	return m.err
}

func (m *mockChatRepo) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) {
	return nil, nil
}

type mockViolationRepo struct{}

func (m *mockViolationRepo) Create(ctx context.Context, violation *models.UserViolation) error {
	return nil
}
func (m *mockViolationRepo) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) {
	return 0, nil
}
func (m *mockViolationRepo) BanUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}

type mockNSFWService struct {
	isNSFW bool
	err    error
}

func (m *mockNSFWService) CheckNSFW(ctx context.Context, file *multipart.FileHeader) (bool, float64, error) {
	return m.isNSFW, 0.0, m.err
}

func TestLocketService_SendLocket(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	var r2Client *storage.R2Client

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, r2Client)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLocketService_SendLocket_MatchNotFound(t *testing.T) {
	matchRepo := &mockMatchRepo{err: errors.New("not found")}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLocketService_SendLocket_CreateError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{err: errors.New("db error")}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLocketService_SendLocket_NSFWDetected(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: true}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "nsfw.jpg"})
	if err == nil {
		t.Errorf("Expected error for NSFW content, got nil")
	}
}

type mockViolationRepo3Strikes struct {
	mockViolationRepo
}

func (m *mockViolationRepo3Strikes) CountActiveViolationsByType(ctx context.Context, userID uuid.UUID, vType string) (int64, error) {
	return 3, nil
}

func TestLocketService_SendLocket_NSFWDetected_3Strikes(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo3Strikes{}
	nsfwService := &mockNSFWService{isNSFW: true}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "nsfw.jpg"})
	if err == nil {
		t.Errorf("Expected error for NSFW content, got nil")
	}
	if err.Error() != "tài khoản của bạn đã bị khóa do vi phạm gửi ảnh nhạy cảm 3 lần" {
		t.Errorf("Expected ban message, got %v", err.Error())
	}
}
