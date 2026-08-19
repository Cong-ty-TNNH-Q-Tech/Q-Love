// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';

class VibeApiClient {
  final Dio dio;

  VibeApiClient({required this.dio});

  Future<Map<String, dynamic>> getCurrentTrack() async {
    try {
      final response = await dio.get('/vibe/current-track');
      return response.data['data'];
    } catch (e) {
      throw Exception('Failed to get current vibe track: $e');
    }
  }

  Future<void> submitMatch(String trackId) async {
    try {
      await dio.post('/vibe/match', data: {
        'track_id': trackId,
      });
    } catch (e) {
      throw Exception('Failed to submit vibe match: $e');
    }
  }
}
