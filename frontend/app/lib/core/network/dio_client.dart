// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:dio/dio.dart';
import 'package:qlove/core/network/token_interceptor.dart';
import 'package:qlove/core/network/secure_storage_service.dart';

class DioClient {
  late final Dio dio;
  final SecureStorageService secureStorageService;

  // We can inject a callback to get the current access token from memory (e.g. from AuthBloc)
  final String? Function() getAccessToken;
  
  // Callback when token is refreshed
  final void Function(String) onTokenRefreshed;
  
  // Callback when refresh token fails and user needs to be logged out
  final void Function() onLogoutRequired;

  DioClient({
    required this.secureStorageService,
    required this.getAccessToken,
    required this.onTokenRefreshed,
    required this.onLogoutRequired,
  }) {
    dio = Dio(BaseOptions(
      baseUrl: const String.fromEnvironment('API_BASE_URL', defaultValue: 'http://10.0.2.2:8080/v1'),
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
    ));

    dio.interceptors.add(TokenInterceptor(
      dio: dio,
      secureStorageService: secureStorageService,
      getAccessToken: getAccessToken,
      onTokenRefreshed: onTokenRefreshed,
      onLogoutRequired: onLogoutRequired,
    ));
  }
}
