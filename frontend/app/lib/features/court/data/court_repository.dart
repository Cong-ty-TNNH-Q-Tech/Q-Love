// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/core/models/court_case_model.dart';

class CourtRepository {
  final Dio dio;

  CourtRepository({required this.dio});

  Future<List<CourtCaseModel>> getCases({int limit = 10, int offset = 0}) async {
    try {
      final response = await dio.get(
        '/court/cases',
        queryParameters: {
          'limit': limit,
          'offset': offset,
        },
      );
      
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data['data'] ?? [];
        return data.map((json) => CourtCaseModel.fromJson(json)).toList();
      }
      throw Exception('Failed to load court cases');
    } catch (e) {
      throw Exception('Error fetching court cases: $e');
    }
  }

  Future<void> vote(String caseId, String voteType) async {
    try {
      final response = await dio.post(
        '/court/cases/$caseId/vote',
        data: {
          'vote': voteType,
        },
      );

      if (response.statusCode != 200) {
        throw Exception('Failed to submit vote');
      }
    } catch (e) {
      throw Exception('Error voting: $e');
    }
  }
}
