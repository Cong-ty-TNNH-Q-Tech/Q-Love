// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/core/network/secure_storage_service.dart';
import 'package:qlove/features/auth/data/auth_repository.dart';

class MockDio extends Mock implements Dio {}
class MockSecureStorageService extends Mock implements SecureStorageService {}

void main() {
  late AuthRepository authRepository;
  late MockDio mockDio;
  late MockSecureStorageService mockSecureStorageService;

  setUp(() {
    mockDio = MockDio();
    mockSecureStorageService = MockSecureStorageService();
    authRepository = AuthRepository(
      dio: mockDio,
      secureStorageService: mockSecureStorageService,
    );
  });

  group('AuthRepository', () {
    test('sendOtp success', () async {
      when(() => mockDio.post('/auth/send-otp', data: {'phone': '0901234567'}))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/auth/send-otp'),
                statusCode: 200,
              ));

      expect(() => authRepository.sendOtp('0901234567'), returnsNormally);
    });

    test('verifyOtp success', () async {
      final mockResponse = {
        'access_token': 'dummy_access_token',
        'refresh_token': 'dummy_refresh_token',
        'is_new_user': true,
        'user': {
          'id': '123',
          'phone': '0901234567',
        }
      };

      when(() => mockDio.post('/auth/verify-otp', data: {'phone': '0901234567', 'otp': '123456'}))
          .thenAnswer((_) async => Response(
                requestOptions: RequestOptions(path: '/auth/verify-otp'),
                statusCode: 200,
                data: mockResponse,
              ));

      when(() => mockSecureStorageService.saveRefreshToken('dummy_refresh_token'))
          .thenAnswer((_) async {});

      final result = await authRepository.verifyOtp('0901234567', '123456');

      expect(result['access_token'], 'dummy_access_token');
      expect(result['is_new_user'], true);
      expect(result['user'], isA<UserModel>());
      verify(() => mockSecureStorageService.saveRefreshToken('dummy_refresh_token')).called(1);
    });
  });
}
