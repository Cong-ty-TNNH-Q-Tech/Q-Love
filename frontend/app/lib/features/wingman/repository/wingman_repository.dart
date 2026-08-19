import 'package:dio/dio.dart';

class WingmanRepository {
  final Dio _dio;

  WingmanRepository(this._dio);

  Future<Map<String, dynamic>> getDashboard() async {
    final response = await _dio.get('/api/v1/wingman/dashboard');
    return response.data as Map<String, dynamic>;
  }

  Future<void> matchFriend(String targetId, String friendId) async {
    await _dio.post('/api/v1/wingman/match', data: {
      'target_id': targetId,
      'friend_id': friendId,
    });
  }
}
