// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	go_redis "github.com/redis/go-redis/v9"
)

// Mocks

type mockUserRepoForAuth struct {
	mock.Mock
}

func (m *mockUserRepoForAuth) GetTopUsersByScore(ctx context.Context, limit int) ([]uuid.UUID, error) {
	return nil, nil
}

func (m *mockUserRepoForAuth) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepoForAuth) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepoForAuth) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return nil, nil
}

type mockESMSClient struct {
	mock.Mock
}

func (m *mockESMSClient) SendOTP(ctx context.Context, phone, otp string) error {
	args := m.Called(ctx, phone, otp)
	return args.Error(0)
}

func setupAuthServiceTest() (*miniredis.Miniredis, *go_redis.Client) {
	mr, _ := miniredis.Run()
	client := go_redis.NewClient(&go_redis.Options{
		Addr: mr.Addr(),
	})
	return mr, client
}

func TestAuthService_SendOTP_Success(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)

	esmsClient.On("SendOTP", mock.Anything, "0901234567", mock.AnythingOfType("string")).Return(nil)

	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	err := svc.SendOTP(context.Background(), "0901234567")
	assert.NoError(t, err)

	esmsClient.AssertExpectations(t)
	
	// Check rate limit and otp cache
	keys, _ := rdb.Keys(context.Background(), "*").Result()
	assert.Len(t, keys, 2) // ratelimit + otp keys
}

func TestAuthService_SendOTP_RateLimitExceeded(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	rdb.Set(context.Background(), "ratelimit:otp:0901234567", 5, time.Hour)

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)

	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	err := svc.SendOTP(context.Background(), "0901234567")
	assert.Error(t, err)
	assert.Equal(t, "rate limit exceeded: max 5 OTPs per day", err.Error())
}

func TestAuthService_VerifyOTP_Success(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	rdb.Set(context.Background(), "otp:0901234567", "123456", 5*time.Minute)

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)

	existingUser := &models.User{
		ID:    uuid.New(),
		Phone: "0901234567",
	}

	userRepo.On("FindByPhone", mock.Anything, "0901234567").Return(existingUser, nil)

	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	user, accessToken, refreshToken, err := svc.VerifyOTP(context.Background(), "0901234567", "123456")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	userRepo.AssertExpectations(t)
}

func TestAuthService_VerifyOTP_CreateNewUser(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	rdb.Set(context.Background(), "otp:0901234567", "123456", 5*time.Minute)

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)

	userRepo.On("FindByPhone", mock.Anything, "0901234567").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	user, accessToken, refreshToken, err := svc.VerifyOTP(context.Background(), "0901234567", "123456")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	userRepo.AssertExpectations(t)
}

func TestAuthService_VerifyOTP_InvalidOTP(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)
	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	// Missing OTP
	_, _, _, err := svc.VerifyOTP(context.Background(), "0901234567", "123456")
	assert.Error(t, err)
	assert.Equal(t, "ERR_INVALID_OTP", err.Error())

	// Wrong OTP
	rdb.Set(context.Background(), "otp:0901234567", "654321", 5*time.Minute)
	_, _, _, err = svc.VerifyOTP(context.Background(), "0901234567", "123456")
	assert.Error(t, err)
	assert.Equal(t, "ERR_INVALID_OTP", err.Error())
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)
	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	// Generate refresh token manually
	rt, _ := svc.(*authService).generateRefreshToken(uuid.New())

	newAT, err := svc.RefreshToken(context.Background(), rt)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAT)
}

func TestAuthService_RefreshToken_Invalid(t *testing.T) {
	mr, rdb := setupAuthServiceTest()
	defer mr.Close()

	userRepo := new(mockUserRepoForAuth)
	esmsClient := new(mockESMSClient)
	svc := NewAuthService(userRepo, esmsClient, rdb, "secret")

	// Use access token instead of refresh
	at, _ := svc.(*authService).generateAccessToken(uuid.New())

	_, err := svc.RefreshToken(context.Background(), at)
	assert.Error(t, err)
	assert.Equal(t, "invalid token type", err.Error())
}
