// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

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

func (m *mockMatchRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func (m *mockMatchRepo) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.match, nil
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

func (m *mockChatRepo) CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error) {
	return 0, nil
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

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
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

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
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

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
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

	fileHeader := &multipart.FileHeader{Filename: "nsfw.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
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

	fileHeader := &multipart.FileHeader{Filename: "nsfw.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil {
		t.Errorf("Expected error for NSFW content, got nil")
	}
	if err.Error() != "tài khoản của bạn đã bị khóa do vi phạm gửi ảnh nhạy cảm 3 lần" {
		t.Errorf("Expected ban message, got %v", err.Error())
	}
}


func TestLocketService_SendLocket_InvalidSize(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 11 * 1024 * 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil || err.Error() != "file too large, limit is 10MB" {
		t.Errorf("Expected file too large error, got %v", err)
	}
}

func TestLocketService_SendLocket_InvalidType(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	fileHeader := &multipart.FileHeader{Filename: "test.pdf", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "application/pdf")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil || err.Error() != "only jpeg, png, and webp images are supported" {
		t.Errorf("Expected invalid type error, got %v", err)
	}
}

func TestLocketService_SendLocket_NSFWError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{err: errors.New("ai error")}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil)

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil || err.Error() != "failed to check image content" {
		t.Errorf("Expected ai error, got %v", err)
	}
}

func createMultipartFile(t *testing.T) *multipart.FileHeader {
	// Instead of a real form, we can just test if the fallback to r2.qlove.com works.
	// But to test R2Client != nil, we need file.Open() to succeed or fail.
	// To make file.Open fail, we can just use an empty FileHeader.
	// Actually file.Open() panics if content is not properly constructed, or returns error.
	return &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
}

func TestLocketService_SendLocket_R2Error(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	// For R2Client to be tested without file.Open() panicking, we rely on file.Open() returning an error
	// because it's a dummy file header!
	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		// S3Client: &mockS3API{err: errors.New("s3 err")}, // Don't even need this if file.Open fails!
	}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, r2Client)

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil {
		t.Errorf("Expected error from file.Open or R2 upload, got nil")
	}
}

func createValidMultipartFileHeader(t *testing.T) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("test image content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_ = req.ParseMultipartForm(1024)
	
	file := req.MultipartForm.File["image"][0]
	file.Header.Set("Content-Type", "image/jpeg")
	return file
}

func TestLocketService_SendLocket_R2UploadError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	// Fake R2 Client with dummy S3Client that will fail network request
	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		S3Client:   &mockS3Client{err: errors.New("s3 upload failed")},
	}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, r2Client)

	fileHeader := createValidMultipartFileHeader(t)
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil {
		t.Errorf("Expected error from R2 upload, got nil")
	}
}

type mockS3Client struct {
	output *s3.PutObjectOutput
	err    error
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.output, m.err
}

func TestLocketService_SendLocket_R2UploadSuccess(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		S3Client:   &mockS3Client{output: &s3.PutObjectOutput{}},
	}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, r2Client)

	fileHeader := createValidMultipartFileHeader(t)
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
}

