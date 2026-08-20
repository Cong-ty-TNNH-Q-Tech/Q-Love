// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockExRatingRepo struct {
	hasRated   bool
	createErr  error
	avg        float64
	total      int64
	tags       map[string]int
	summaryErr error
}

func (m *mockExRatingRepo) Create(ctx context.Context, rating *models.ExRating) error {
	return m.createErr
}
func (m *mockExRatingRepo) HasRated(ctx context.Context, matchID, targetUserID uuid.UUID) (bool, error) {
	return m.hasRated, nil
}
func (m *mockExRatingRepo) GetSummaryByUserID(ctx context.Context, targetUserID uuid.UUID) (float64, int64, map[string]int, error) {
	return m.avg, m.total, m.tags, m.summaryErr
}

type mockChatRepoForExRating struct {
	msgCount int64
}

func (m *mockChatRepoForExRating) Create(ctx context.Context, msg *models.ChatMessage) error { return nil }
func (m *mockChatRepoForExRating) GetMessagesByMatchID(ctx context.Context, matchID uuid.UUID, limit int, before *time.Time) ([]models.ChatMessage, error) { return nil, nil }
func (m *mockChatRepoForExRating) CountMessagesByMatchID(ctx context.Context, matchID uuid.UUID) (int64, error) {
	return m.msgCount, nil
}

type mockMatchRepoForExRating struct {
	match *models.Match
}

func (m *mockMatchRepoForExRating) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	if m.match == nil {
		return nil, errors.New("not found")
	}
	return m.match, nil
}
func (m *mockMatchRepoForExRating) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockMatchRepoForExRating) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	if m.match == nil {
		return nil, errors.New("not found")
	}
	return m.match, nil
}
func (m *mockMatchRepoForExRating) Create(ctx context.Context, match *models.Match) error { return nil }
func (m *mockMatchRepoForExRating) UpdateLastInteraction(ctx context.Context, matchID uuid.UUID, t time.Time) error { return nil }

type mockWalletRepoForExRating struct {
	balance float64
}
func (m *mockWalletRepoForExRating) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (*models.UserWallet, error) {
	return &models.UserWallet{Balance: m.balance}, nil
}
func (m *mockWalletRepoForExRating) UpdateBalance(ctx context.Context, userID uuid.UUID, delta float64) error {
	return nil
}
func (m *mockWalletRepoForExRating) AddCommission(ctx context.Context, userID uuid.UUID, amount float64) error {
	return nil
}
func (m *mockWalletRepoForExRating) CreateTransaction(ctx context.Context, txn *models.WalletTransaction) error {
	return nil
}
func (m *mockWalletRepoForExRating) CheckTransactionExists(ctx context.Context, txID uuid.UUID) (bool, error) {
	return false, nil
}

type mockTxManagerForExRating struct{}
func (m *mockTxManagerForExRating) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error, opts ...*sql.TxOptions) error {
	return fn(ctx)
}

func TestExRatingService_SubmitRating_Success(t *testing.T) {
	targetUserID := uuid.New()
	matchID := uuid.New()
	match := &models.Match{ID: matchID, User1ID: uuid.New(), User2ID: targetUserID, DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}

	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: false},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 51},
		&mockMatchRepoForExRating{match: match},
	)

	err := svc.SubmitRating(context.Background(), targetUserID, matchID, 4, []string{"#tốt"})
	assert.NoError(t, err)
}

func TestExRatingService_SubmitRating_NotUnmatched(t *testing.T) {
	targetUserID := uuid.New()
	matchID := uuid.New()
	match := &models.Match{ID: matchID, User1ID: uuid.New(), User2ID: targetUserID, DeletedAt: gorm.DeletedAt{Valid: false}}

	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: false},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 51},
		&mockMatchRepoForExRating{match: match},
	)

	err := svc.SubmitRating(context.Background(), targetUserID, matchID, 4, []string{"#tốt"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmatch")
}

func TestExRatingService_SubmitRating_NotEnoughMessages(t *testing.T) {
	targetUserID := uuid.New()
	matchID := uuid.New()
	match := &models.Match{ID: matchID, User1ID: uuid.New(), User2ID: targetUserID, DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}

	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: false},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 49}, // < 50
		&mockMatchRepoForExRating{match: match},
	)

	err := svc.SubmitRating(context.Background(), targetUserID, matchID, 4, []string{"#tốt"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "50")
}

func TestExRatingService_SubmitRating_MatchNotFound(t *testing.T) {
	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: false},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 51},
		&mockMatchRepoForExRating{match: nil},
	)

	err := svc.SubmitRating(context.Background(), uuid.New(), uuid.New(), 4, []string{"#tốt"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestExRatingService_SubmitRating_TargetUserNotMatch(t *testing.T) {
	match := &models.Match{ID: uuid.New(), User1ID: uuid.New(), User2ID: uuid.New(), DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}
	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: false},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 51},
		&mockMatchRepoForExRating{match: match},
	)

	err := svc.SubmitRating(context.Background(), uuid.New(), match.ID, 4, []string{"#tốt"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "thuộc về")
}

func TestExRatingService_SubmitRating_AlreadyRated(t *testing.T) {
	targetUserID := uuid.New()
	matchID := uuid.New()
	match := &models.Match{ID: matchID, User1ID: uuid.New(), User2ID: targetUserID, DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}
	svc := NewExRatingService(
		&mockExRatingRepo{hasRated: true},
		&mockWalletRepoForExRating{balance: 100},
		&mockTxManagerForExRating{},
		&mockChatRepoForExRating{msgCount: 51},
		&mockMatchRepoForExRating{match: match},
	)

	err := svc.SubmitRating(context.Background(), targetUserID, matchID, 4, []string{"#tốt"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "đã đánh giá")
}

func TestExRatingService_ViewRating_Success(t *testing.T) {
	viewerID := uuid.New()
	targetID := uuid.New()

	svc := NewExRatingService(
		&mockExRatingRepo{avg: 4.5, total: 10, tags: map[string]int{"#green": 5}},
		&mockWalletRepoForExRating{balance: 100}, // > 50 xu
		&mockTxManagerForExRating{},
		nil,
		nil,
	)

	avg, total, tags, err := svc.ViewRating(context.Background(), viewerID, targetID)
	assert.NoError(t, err)
	assert.Equal(t, 4.5, avg)
	assert.Equal(t, int64(10), total)
	assert.Equal(t, 5, tags["#green"])
}

func TestExRatingService_ViewRating_InsufficientFunds(t *testing.T) {
	viewerID := uuid.New()
	targetID := uuid.New()

	svc := NewExRatingService(
		&mockExRatingRepo{},
		&mockWalletRepoForExRating{balance: 49}, // < 50 xu
		&mockTxManagerForExRating{},
		nil,
		nil,
	)

	_, _, _, err := svc.ViewRating(context.Background(), viewerID, targetID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "xu")
}

type mockTxManagerWithError struct{}

func (m *mockTxManagerWithError) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return errors.New("lỗi truy xuất ví")
}

func TestExRatingService_ViewRating_WalletError(t *testing.T) {
	viewerID := uuid.New()
	targetID := uuid.New()

	svc := NewExRatingService(
		&mockExRatingRepo{},
		&mockWalletRepoForExRating{balance: 100}, 
		&mockTxManagerWithError{},
		nil,
		nil,
	)

	_, _, _, err := svc.ViewRating(context.Background(), viewerID, targetID)
	assert.Error(t, err)
}
