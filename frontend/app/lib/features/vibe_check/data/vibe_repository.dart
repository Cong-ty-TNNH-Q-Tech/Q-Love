// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'vibe_api_client.dart';

class VibeRepository {
  final VibeApiClient apiClient;

  VibeRepository({required this.apiClient});

  Future<Map<String, String>> getCurrentVibeTrack() async {
    final trackData = await apiClient.getCurrentTrack();
    return {
      "id": trackData['id'] ?? "",
      "title": trackData['name'] ?? "Unknown Title",
      "artist": trackData['artist'] ?? "Unknown Artist",
      "coverUrl": trackData['album_art'] ?? "https://example.com/default_art.jpg",
      "previewUrl": "", // Backend does not return preview URL currently
    };
  }

  Future<void> matchVibe(String trackId) async {
    await apiClient.submitMatch(trackId);
  }
}
