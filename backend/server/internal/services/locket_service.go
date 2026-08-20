package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/imageutil"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/storage"
	"github.com/google/uuid"
)

type LocketService interface {
	SendLocket(ctx context.Context, senderID uuid.UUID, matchID uuid.UUID, file *multipart.FileHeader) error
}

type locketService struct {
	chatRepo      repository.ChatRepository
	matchRepo     repository.MatchRepository
	violationRepo repository.UserViolationRepository
	nsfwService   NSFWService
	r2Client      *storage.R2Client
}

func NewLocketService(
	chatRepo repository.ChatRepository,
	matchRepo repository.MatchRepository,
	violationRepo repository.UserViolationRepository,
	nsfwService NSFWService,
	r2Client *storage.R2Client,
) LocketService {
	return &locketService{
		chatRepo:      chatRepo,
		matchRepo:     matchRepo,
		violationRepo: violationRepo,
		nsfwService:   nsfwService,
		r2Client:      r2Client,
	}
}

func (s *locketService) SendLocket(ctx context.Context, senderID uuid.UUID, matchID uuid.UUID, file *multipart.FileHeader) error {
	// Verify match exists
	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return errors.New("match not found")
	}

	// AI Check NSFW
	isNSFW, _, err := s.nsfwService.CheckNSFW(ctx, file)
	if err != nil {
		return errors.New("failed to check image content")
	}

	if isNSFW {
		// 1. Log violation
		violation := &models.UserViolation{
			UserID: senderID,
			Type:   "nsfw_image",
			Reason: "System detected high skin ratio (>30%)",
		}
		_ = s.violationRepo.Create(ctx, violation)

		// 2. Check 3 strikes rule
		count, _ := s.violationRepo.CountActiveViolationsByType(ctx, senderID, "nsfw_image")
		if count >= 3 {
			_ = s.violationRepo.BanUser(ctx, senderID)
			return errors.New("tài khoản của bạn đã bị khóa do vi phạm gửi ảnh nhạy cảm 3 lần")
		}

		return errors.New("ảnh chứa nội dung nhạy cảm, không được phép gửi")
	}

	// Validate file type and size
	if file.Size > 10*1024*1024 { // 10MB limit
		return errors.New("file too large, limit is 10MB")
	}
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return errors.New("only jpeg, png, and webp images are supported")
	}

	// 3. Process blur
	// Default max streak for 0% blur is 30 days
	blurPercentage := 100
	if match.StreakScore > 0 {
		blurPercentage = 100 - int(match.StreakScore*100/30)
		if blurPercentage < 0 {
			blurPercentage = 0
		}
	}

	// Read raw image
	var srcImage image.Image
	f, openErr := file.Open()
	if openErr != nil {
		return errors.New("failed to read image")
	}
	defer f.Close()
	
	srcImage, _, err = image.Decode(f)
	if err != nil {
		return errors.New("failed to decode image")
	}

	// Apply Gaussian Blur approximation
	blurredImage := imageutil.ApplyGaussianBlur(srcImage, blurPercentage)

	// Encode back to buffer
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, blurredImage, &jpeg.Options{Quality: 85}); err != nil {
		return errors.New("failed to process image")
	}

	// Upload using R2Client
	var imageURL string
	if s.r2Client != nil {
		objectKey := fmt.Sprintf("lockets/%s/%s/blurred_%s", matchID.String(), senderID.String(), file.Filename)
		url, err := s.r2Client.UploadFile(ctx, objectKey, buf, "image/jpeg")
		if err != nil {
			return fmt.Errorf("failed to upload blurred image to R2: %w", err)
		}
		imageURL = url
	} else {
		imageURL = "https://r2.qlove.com/" + matchID.String() + "/" + file.Filename
	}

	blurURL := imageURL // The image itself is blurred

	chatMessage := &models.ChatMessage{
		MatchID:  matchID,
		SenderID: senderID,
		Type:     "locket",
		Content:  imageURL,
		BlurURL:  blurURL,
	}

	if err := s.chatRepo.Create(ctx, chatMessage); err != nil {
		return err
	}

	// Update last interaction
	_ = s.matchRepo.UpdateLastInteraction(ctx, matchID, chatMessage.CreatedAt)

	return nil
}
