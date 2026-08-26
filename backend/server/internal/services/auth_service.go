// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/esms"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	SendOTP(ctx context.Context, phone string) error
	VerifyOTP(ctx context.Context, phone, otp string) (*models.User, string, string, bool, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	esmsClient esms.Client
	redis      *redis.Client
	jwtSecret  string
}

func NewAuthService(userRepo repository.UserRepository, esmsClient esms.Client, redis *redis.Client, jwtSecret string) AuthService {
	return &authService{
		userRepo:   userRepo,
		esmsClient: esmsClient,
		redis:      redis,
		jwtSecret:  jwtSecret,
	}
}

// Generate 6-digit OTP
func generateOTP() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *authService) SendOTP(ctx context.Context, phone string) error {
	// Rate limit: 5 times per day per phone
	rateKey := fmt.Sprintf("ratelimit:otp:%s", phone)
	count, err := s.redis.Get(ctx, rateKey).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if count >= 5 {
		return errors.New("rate limit exceeded: max 5 OTPs per day")
	}

	// Generate OTP
	otp := generateOTP()
	
	// Cache OTP for 120s
	otpKey := fmt.Sprintf("otp:%s", phone)
	err = s.redis.Set(ctx, otpKey, otp, 120*time.Second).Err()
	if err != nil {
		return err
	}

	// Send OTP using ESMS client
	err = s.esmsClient.SendOTP(ctx, phone, otp)
	if err != nil {
		return err
	}

	// Increment rate limit
	pipe := s.redis.Pipeline()
	pipe.Incr(ctx, rateKey)
	if count == 0 {
		pipe.Expire(ctx, rateKey, 24*time.Hour)
	}
	_, err = pipe.Exec(ctx)
	
	return err
}

func (s *authService) VerifyOTP(ctx context.Context, phone, otp string) (*models.User, string, string, bool, error) {
	otpKey := fmt.Sprintf("otp:%s", phone)
	attemptsKey := fmt.Sprintf("otp_attempts:%s", phone)

	cachedOTP, err := s.redis.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return nil, "", "", false, errors.New("ERR_INVALID_OTP")
	} else if err != nil {
		return nil, "", "", false, err
	}

	if cachedOTP != otp {
		attempts, _ := s.redis.Incr(ctx, attemptsKey).Result()
		s.redis.Expire(ctx, attemptsKey, 120*time.Second)
		if attempts >= 3 {
			s.redis.Del(ctx, otpKey)
			s.redis.Del(ctx, attemptsKey)
			return nil, "", "", false, errors.New("ERR_TOO_MANY_ATTEMPTS")
		}
		return nil, "", "", false, errors.New("ERR_INVALID_OTP")
	}

	// Remove OTP after successful verification
	s.redis.Del(ctx, otpKey)
	s.redis.Del(ctx, attemptsKey)

	// Check if user exists
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, "", "", false, err
	}

	isNewUser := false
	if user == nil {
		isNewUser = true
		// Create new user
		user = &models.User{
			Phone: phone,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, "", "", false, err
		}
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", false, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", false, err
	}

	return user, accessToken, refreshToken, isNewUser, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	if typ, ok := claims["type"].(string); !ok || typ != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("missing sub in claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
