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

type mockChatRepo struct {
	err error
}

func (m *mockChatRepo) Create(ctx context.Context, message *models.ChatMessage) error {
	return m.err
}

func TestLocketService_SendLocket(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	
	// Just passing nil for r2Client as we are mocking inside the service for now 
	// (the real logic is simulated in the service).
	var r2Client *storage.R2Client
	
	service := NewLocketService(chatRepo, matchRepo, r2Client)
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLocketService_SendLocket_MatchNotFound(t *testing.T) {
	matchRepo := &mockMatchRepo{err: errors.New("not found")}
	chatRepo := &mockChatRepo{}
	
	service := NewLocketService(chatRepo, matchRepo, nil)
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLocketService_SendLocket_CreateError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{err: errors.New("db error")}
	
	service := NewLocketService(chatRepo, matchRepo, nil)
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), &multipart.FileHeader{Filename: "test.jpg"})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
