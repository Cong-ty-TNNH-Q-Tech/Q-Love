// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/core/models/user_model.dart';

class DiscoverRepository {
  final Dio dio;

  DiscoverRepository({required this.dio});

  Future<List<UserModel>> getFeed({String filter = 'default', int limit = 10}) async {
    try {
      final response = await dio.get(
        '/users/feed',
        queryParameters: {
          'filter': filter,
          'limit': limit,
        },
      );
      
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data;
        return data.map((json) => UserModel.fromJson(json)).toList();
      }
      throw Exception('Failed to load feed');
    } catch (e) {
      throw Exception('Error fetching feed: $e');
    }
  }

  Future<bool> swipe(String targetId, String action) async {
    try {
      final response = await dio.post(
        '/matches/swipe',
        data: {
          'target_id': targetId,
          'action': action,
        },
      );

      if (response.statusCode == 200) {
        return response.data['is_match'] == true;
      }
      return false;
    } catch (e) {
      throw Exception('Error swiping: $e');
    }
  }
}
