// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (m *mockMatchRepo) FindByUsers(ctx context.Context, u1, u2 uuid.UUID) (*models.Match, error) {
	return nil, nil
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

func (m *mockViolationRepo) GetViolations(ctx context.Context, page, limit int) ([]models.UserViolation, int64, error) {
	return nil, 0, nil
}

func (m *mockViolationRepo) DeleteViolation(ctx context.Context, violationID uuid.UUID) error {
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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), createValidMultipartFile(t))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLocketService_SendLocket_MatchNotFound(t *testing.T) {
	matchRepo := &mockMatchRepo{err: errors.New("not found")}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), createValidMultipartFile(t))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLocketService_SendLocket_NSFWDetected(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: true}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

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

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

	fileHeader := &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil || err.Error() != "failed to check image content" {
		t.Errorf("Expected ai error, got %v", err)
	}
}

func createMultipartFile(t *testing.T) *multipart.FileHeader {
	return &multipart.FileHeader{Filename: "test.jpg", Size: 1024, Header: make(map[string][]string)}
}

func createValidMultipartFile(t *testing.T) *multipart.FileHeader {
	// Create a minimal valid JPEG image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(buf.Bytes())
	writer.Close()

	req := &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(body),
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	reader, err := req.MultipartReader()
	if err != nil {
		t.Fatal(err)
	}
	form, err := reader.ReadForm(1024 * 1024)
	if err != nil {
		t.Fatal(err)
	}

	fileHeader := form.File["file"][0]
	if fileHeader.Header == nil {
		fileHeader.Header = make(map[string][]string)
	}
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	return fileHeader
}

func createInvalidMultipartFile(t *testing.T) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("this is not a valid image file"))
	writer.Close()

	req := &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(body),
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	reader, err := req.MultipartReader()
	if err != nil {
		t.Fatal(err)
	}
	form, err := reader.ReadForm(1024 * 1024)
	if err != nil {
		t.Fatal(err)
	}
	fileHeader := form.File["file"][0]
	fileHeader.Header.Set("Content-Type", "image/jpeg")
	return fileHeader
}

func TestLocketService_SendLocket_DecodeError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, nil)

	fileHeader := createInvalidMultipartFile(t)
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil || err.Error() != "failed to decode image" {
		t.Errorf("Expected failed to decode image error, got %v", err)
	}
}

func TestLocketService_SendLocket_R2Error(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		S3Client:   s3.New(s3.Options{}),
	}
	// This will cause a network error instead of a panic since S3Client is not nil but unconfigured
	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	fileHeader := createValidMultipartFile(t)
	
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil {
		t.Errorf("Expected error from R2 upload, got nil")
	}
}

func TestLocketService_SendLocket_WithStreak(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{StreakScore: 15}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	var r2Client *storage.R2Client

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), createValidMultipartFile(t))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLocketService_SendLocket_WithHighStreak(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{StreakScore: 40}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	var r2Client *storage.R2Client

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), createValidMultipartFile(t))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

type mockS3Client struct {
	output *s3.PutObjectOutput
	err    error
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.output, m.err
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return &s3.DeleteObjectOutput{}, m.err
}

func TestLocketService_SendLocket_R2UploadSuccess(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{StreakScore: 0}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		S3Client:   &mockS3Client{output: &s3.PutObjectOutput{}},
	}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	fileHeader := createValidMultipartFile(t)
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}
}

func TestLocketService_SendLocket_R2UploadError(t *testing.T) {
	matchRepo := &mockMatchRepo{match: &models.Match{StreakScore: 0}}
	chatRepo := &mockChatRepo{}
	violationRepo := &mockViolationRepo{}
	nsfwService := &mockNSFWService{isNSFW: false}

	r2Client := &storage.R2Client{
		BucketName: "test-bucket",
		S3Client:   &mockS3Client{err: errors.New("s3 upload failed")},
	}

	service := NewLocketService(chatRepo, matchRepo, violationRepo, nsfwService, nil, r2Client)

	fileHeader := createValidMultipartFile(t)
	err := service.SendLocket(context.Background(), uuid.New(), uuid.New(), fileHeader)
	if err == nil {
		t.Errorf("Expected error from R2 upload, got nil")
	}
}

