// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"math/rand"
	"time"
)

type SpotifyService struct {
	// dependencies like Redis client for token caching would go here
}

func NewSpotifyService() *SpotifyService {
	return &SpotifyService{}
}

type Track struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	AlbumArt string `json:"album_art"`
}

var mockTracks = []Track{
	{ID: "track1", Name: "Shape of You", Artist: "Ed Sheeran", AlbumArt: "https://example.com/art1.jpg"},
	{ID: "track2", Name: "Blinding Lights", Artist: "The Weeknd", AlbumArt: "https://example.com/art2.jpg"},
	{ID: "track3", Name: "Levitating", Artist: "Dua Lipa", AlbumArt: "https://example.com/art3.jpg"},
}

// GetCurrentTrack mocks fetching the user's currently playing track from Spotify
func (s *SpotifyService) GetCurrentTrack(userID string) (Track, error) {
	// In reality, this would use the user's cached Spotify access token to call the Spotify API
	// Mock: return a random track
	rand.Seed(time.Now().UnixNano())
	track := mockTracks[rand.Intn(len(mockTracks))]
	return track, nil
}

var TimeNow = time.Now

// CheckUnlockTime returns whether it is >= 23:00 and < 05:00
func (s *SpotifyService) CheckUnlockTime() bool {
	now := TimeNow()
	hour := now.Hour()
	if hour >= 23 || hour < 5 {
		return true
	}
	return false
}
