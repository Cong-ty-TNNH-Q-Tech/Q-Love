// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';

class MinigameApiClient {
  final Dio dio;

  MinigameApiClient({required this.dio});

  Future<String> initSteal(String defenderId, String targetCardId) async {
    try {
      final response = await dio.post('/api/v1/minigame/steal/init', data: {
        'defender_id': defenderId,
        'target_card_id': targetCardId,
      });
      return response.data['data']['id'];
    } catch (e) {
      throw Exception('Failed to init steal session: $e');
    }
  }

  Future<void> submitStealResult(String stealId, bool isWin) async {
    try {
      await dio.post('/api/v1/minigame/steal/submit', data: {
        'steal_id': stealId,
        'is_win': isWin,
      });
    } catch (e) {
      throw Exception('Failed to submit steal result: $e');
    }
  }
}
