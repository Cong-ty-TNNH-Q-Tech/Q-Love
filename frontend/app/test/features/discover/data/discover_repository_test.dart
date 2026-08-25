// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/features/discover/data/discover_repository.dart';

class MockDio extends Mock implements Dio {}

void main() {
  group('DiscoverRepository', () {
    late DiscoverRepository repository;
    late MockDio mockDio;

    setUp(() {
      mockDio = MockDio();
      repository = DiscoverRepository(dio: mockDio);
    });

    test('getFeed returns list of UserModel on success', () async {
      when(() => mockDio.get(any(), queryParameters: any(named: 'queryParameters')))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/users/feed'),
                statusCode: 200,
                data: [
                  {'id': '1', 'name': 'Test User 1'},
                  {'id': '2', 'name': 'Test User 2'},
                ],
              ));

      final result = await repository.getFeed();

      expect(result.length, 2);
      expect(result.first.id, '1');
      expect(result.first.name, 'Test User 1');
    });

    test('getFeed throws Exception on failure', () async {
      when(() => mockDio.get(any(), queryParameters: any(named: 'queryParameters')))
          .thenThrow(DioException(requestOptions: RequestOptions(path: '/users/feed')));

      expect(() => repository.getFeed(), throwsException);
    });

    test('swipe calls endpoint and returns match status', () async {
      when(() => mockDio.post(any(), data: any(named: 'data')))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/matches/swipe'),
                statusCode: 200,
                data: {
                  'is_match': true,
                },
              ));

      final result = await repository.swipe('1', 'like');

      expect(result, true);
    });

    test('swipe throws Exception on failure', () async {
      when(() => mockDio.post(any(), data: any(named: 'data')))
          .thenThrow(DioException(requestOptions: RequestOptions(path: '/matches/swipe')));

      expect(() => repository.swipe('1', 'like'), throwsException);
    });
  });
}
