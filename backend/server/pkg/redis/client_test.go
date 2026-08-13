package redis

import (
	"testing"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/stretchr/testify/assert"
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
