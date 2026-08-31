// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

const (
	purgeQueueKey = "qlove:purge:queue"
)

type PurgeQueueRepository interface {
	EnqueueUser(ctx context.Context, userID string, isVIP bool) error
	DequeueUsers(ctx context.Context, count int64) ([]string, error)
	DequeueVIPUsers(ctx context.Context, count int64) ([]string, error)
	ClearQueue(ctx context.Context) error
}

type purgeQueueRepositoryImpl struct {
	rdb *redis.Client
}

func NewPurgeQueueRepository(rdb *redis.Client) PurgeQueueRepository {
	return &purgeQueueRepositoryImpl{rdb: rdb}
}

func (r *purgeQueueRepositoryImpl) EnqueueUser(ctx context.Context, userID string, isVIP bool) error {
	key := purgeQueueKey
	if isVIP {
		key = purgeQueueKey + ":vip"
	}
	return r.rdb.LPush(ctx, key, userID).Err()
}

func (r *purgeQueueRepositoryImpl) DequeueUsers(ctx context.Context, count int64) ([]string, error) {
	res, err := r.rdb.RPopCount(ctx, purgeQueueKey, int(count)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return res, err
}

func (r *purgeQueueRepositoryImpl) DequeueVIPUsers(ctx context.Context, count int64) ([]string, error) {
	res, err := r.rdb.RPopCount(ctx, purgeQueueKey+":vip", int(count)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return res, err
}

func (r *purgeQueueRepositoryImpl) ClearQueue(ctx context.Context) error {
	pipe := r.rdb.Pipeline()
	pipe.Del(ctx, purgeQueueKey)
	pipe.Del(ctx, purgeQueueKey+":vip")
	_, err := pipe.Exec(ctx)
	return err
}
