// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/features/drama/models/shame_model.dart';

abstract class ShameRepository {
  Future<List<ShameModel>> getActiveShames({int limit = 10, int offset = 0});
  Future<void> throwTomato(String shameId);
}

class ShameApiRepository implements ShameRepository {
  final Dio _dio;

  ShameApiRepository({required Dio dio}) : _dio = dio;

  @override
  Future<List<ShameModel>> getActiveShames({int limit = 10, int offset = 0}) async {
    final response = await _dio.get('/shames/active', queryParameters: {
      'limit': limit,
      'offset': offset,
    });

    final data = response.data['data'] as List;
    return data.map((json) => ShameModel.fromJson(json)).toList();
  }

  @override
  Future<void> throwTomato(String shameId) async {
    await _dio.post('/shames/$shameId/tomato');
  }
}
