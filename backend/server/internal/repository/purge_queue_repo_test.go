// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPurgeQueueRepository_EnqueueUser(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewPurgeQueueRepository(db)

	userID := "user123"

	// Success case
	mock.ExpectLPush("qlove:purge:queue", userID).SetVal(1)
	err := repo.EnqueueUser(context.Background(), userID, false)
	assert.NoError(t, err)

	// VIP case
	mock.ExpectLPush("qlove:purge:queue:vip", userID).SetVal(1)
	err = repo.EnqueueUser(context.Background(), userID, true)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPurgeQueueRepository_DequeueUsers(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewPurgeQueueRepository(db)

	// Success case normal
	mock.ExpectRPopCount("qlove:purge:queue", 2).SetVal([]string{"user1", "user2"})
	users, err := repo.DequeueUsers(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Success case VIP
	mock.ExpectRPopCount("qlove:purge:queue:vip", 2).SetVal([]string{"vip1", "vip2"})
	vipUsers, err := repo.DequeueVIPUsers(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, vipUsers, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}



func TestPurgeQueueRepository_ClearQueue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := NewPurgeQueueRepository(db)

	mock.ExpectDel("qlove:purge:queue", "qlove:purge:queue:vip").SetVal(2)
	err := repo.ClearQueue(context.Background())
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
