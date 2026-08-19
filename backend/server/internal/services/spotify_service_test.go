// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpotifyService_GetCurrentTrack(t *testing.T) {
	service := NewSpotifyService()
	track, err := service.GetCurrentTrack("user1")
	assert.NoError(t, err)
	assert.NotEmpty(t, track.ID)
	assert.NotEmpty(t, track.Name)
}

func TestSpotifyService_CheckUnlockTime(t *testing.T) {
	service := NewSpotifyService()
	unlocked := service.CheckUnlockTime()
	
	// Check against the current actual time
	hour := time.Now().Hour()
	expected := (hour >= 23 || hour < 5)
	assert.Equal(t, expected, unlocked)
}
