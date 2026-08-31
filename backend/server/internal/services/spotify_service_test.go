// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

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
	assert.NotEmpty(t, track.Artist)
	assert.NotEmpty(t, track.AlbumArt)
}

func TestSpotifyService_GetCurrentTrack_ReturnsValidTrack(t *testing.T) {
	service := NewSpotifyService()
	validIDs := map[string]bool{"track1": true, "track2": true, "track3": true}

	for i := 0; i < 10; i++ {
		track, err := service.GetCurrentTrack("user1")
		assert.NoError(t, err)
		assert.True(t, validIDs[track.ID], "track ID should be one of the mock tracks")
	}
}

func TestSpotifyService_GetCurrentTrack_DifferentUsers(t *testing.T) {
	service := NewSpotifyService()
	_, err1 := service.GetCurrentTrack("user1")
	_, err2 := service.GetCurrentTrack("user2")
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestSpotifyService_CheckUnlockTime(t *testing.T) {
	service := NewSpotifyService()

	tests := []struct {
		name     string
		hour     int
		expected bool
	}{
		{"23:00 - unlocked", 23, true},
		{"00:00 - unlocked", 0, true},
		{"01:00 - unlocked", 1, true},
		{"04:00 - unlocked", 4, true},
		{"05:00 - locked", 5, false},
		{"12:00 - locked", 12, false},
		{"18:00 - locked", 18, false},
		{"22:00 - locked", 22, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TimeNow = func() time.Time {
				return time.Date(2026, 1, 1, tt.hour, 0, 0, 0, time.UTC)
			}
			result := service.CheckUnlockTime()
			assert.Equal(t, tt.expected, result)
		})
	}

	// Restore TimeNow
	TimeNow = time.Now
}

func TestNewSpotifyService(t *testing.T) {
	service := NewSpotifyService()
	assert.NotNil(t, service)
}

func TestMockTracks_NotEmpty(t *testing.T) {
	assert.NotEmpty(t, mockTracks)
	for _, track := range mockTracks {
		assert.NotEmpty(t, track.ID)
		assert.NotEmpty(t, track.Name)
		assert.NotEmpty(t, track.Artist)
		assert.NotEmpty(t, track.AlbumArt)
	}
}
