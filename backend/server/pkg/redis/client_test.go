// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package redis

import (
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/alicebob/miniredis/v2"
)

func TestNewRedisClient_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		RedisURL: "invalid-url",
	}

	client, err := NewRedisClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewRedisClient_PingFail(t *testing.T) {
	cfg := &config.Config{
		RedisURL: "redis://localhost:9999/0", // assuming no redis runs on 9999
	}

	client, err := NewRedisClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewRedisClient_Skip(t *testing.T) {
	cfg := &config.Config{
		RedisURL: "skip",
	}
	client, err := NewRedisClient(cfg)
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestNewRedisClient_Success(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	cfg := &config.Config{
		RedisURL: "redis://" + mr.Addr(),
	}

	client, err := NewRedisClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
