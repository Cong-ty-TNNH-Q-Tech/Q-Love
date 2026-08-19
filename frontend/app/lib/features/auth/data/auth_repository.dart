// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/core/network/secure_storage_service.dart';

class AuthRepository {
  final Dio _dio;
  final SecureStorageService _secureStorageService;

  AuthRepository({
    required Dio dio,
    required SecureStorageService secureStorageService,
  })  : _dio = dio,
        _secureStorageService = secureStorageService;

  Future<void> sendOtp(String phone) async {
    try {
      await _dio.post('/auth/send-otp', data: {'phone': phone});
    } on DioException catch (e) {
      throw Exception(e.response?.data['message'] ?? 'Failed to send OTP');
    } catch (e) {
      throw Exception('An unexpected error occurred');
    }
  }

  Future<Map<String, dynamic>> verifyOtp(String phone, String otp) async {
    try {
      final response = await _dio.post(
        '/auth/verify-otp',
        data: {'phone': phone, 'otp': otp},
      );

      final accessToken = response.data['access_token'] as String;
      final refreshToken = response.data['refresh_token'] as String;
      final isNewUser = response.data['is_new_user'] as bool;
      final user = UserModel.fromJson(response.data['user']);

      // Save refresh token
      await _secureStorageService.saveRefreshToken(refreshToken);

      return {
        'access_token': accessToken,
        'is_new_user': isNewUser,
        'user': user,
      };
    } on DioException catch (e) {
      throw Exception(e.response?.data['message'] ?? 'Failed to verify OTP');
    } catch (e) {
      throw Exception('An unexpected error occurred');
    }
  }

  Future<UserModel> updateProfile({
    required String name,
    required String gender,
    String? avatarUrl,
    String? zodiac,
  }) async {
    try {
      // Dựa theo API.yaml, PUT /users/me dùng để cập nhật profile.
      final response = await _dio.put(
        '/users/me',
        data: {
          'name': name,
          'gender': gender,
          if (avatarUrl != null) 'avatar_url': avatarUrl,
          if (zodiac != null) 'zodiac': zodiac,
        },
      );
      
      return UserModel.fromJson(response.data);
    } on DioException catch (e) {
      throw Exception(e.response?.data['message'] ?? 'Failed to update profile');
    } catch (e) {
      throw Exception('An unexpected error occurred');
    }
  }

  Future<void> logout() async {
    await _secureStorageService.deleteRefreshToken();
  }
}
