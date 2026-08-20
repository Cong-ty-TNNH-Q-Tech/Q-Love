// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package redis

import (
	"context"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	if cfg.RedisURL == "skip" {
		return nil, nil
	}

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	
	client := redis.NewClient(opts)
	
	// Ping to ensure connection is alive
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	
	return client, nil
}
