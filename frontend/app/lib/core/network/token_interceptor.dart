// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/core/network/secure_storage_service.dart';

class TokenInterceptor extends Interceptor {
  final Dio dio;
  final SecureStorageService secureStorageService;
  final String? Function() getAccessToken;
  final void Function(String) onTokenRefreshed;
  final void Function() onLogoutRequired;

  TokenInterceptor({
    required this.dio,
    required this.secureStorageService,
    required this.getAccessToken,
    required this.onTokenRefreshed,
    required this.onLogoutRequired,
  });

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    // Add access token to headers if available
    final accessToken = getAccessToken();
    if (accessToken != null && accessToken.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $accessToken';
    }
    super.onRequest(options, handler);
  }

  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    // If the error is 401 Unauthorized and the path is not refresh token, try to refresh the token
    if (err.response?.statusCode == 401 && err.requestOptions.path != '/auth/refresh') {
      final refreshToken = await secureStorageService.getRefreshToken();
      if (refreshToken != null && refreshToken.isNotEmpty) {
        try {
          // Attempt to refresh
          final refreshResponse = await dio.post(
            '/auth/refresh',
            data: {'refresh_token': refreshToken},
          );

          if (refreshResponse.statusCode == 200) {
            final newAccessToken = refreshResponse.data['access_token'];
            final newRefreshToken = refreshResponse.data['refresh_token'];

            // Save new refresh token
            if (newRefreshToken != null) {
              await secureStorageService.saveRefreshToken(newRefreshToken);
            }

            // Retry the original request with the new access token
            final options = err.requestOptions;
            options.headers['Authorization'] = 'Bearer $newAccessToken';
            
            onTokenRefreshed(newAccessToken);
            
            final retryResponse = await dio.fetch(options);
            return handler.resolve(retryResponse);
          }
        } catch (e) {
          // Refresh failed, logout
          onLogoutRequired();
        }
      } else {
        // No refresh token, logout
        onLogoutRequired();
      }
    }
    
    super.onError(err, handler);
  }
}
