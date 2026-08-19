// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:dio/dio.dart';
import 'chat_model.dart';

class ChatRepository {
  WebSocketChannel? _channel;
  final Dio _dio = Dio();
  final String baseUrl = const String.fromEnvironment('API_URL', defaultValue: 'http://10.0.2.2:8080/api/v1');
  final String wsUrl = const String.fromEnvironment('WS_URL', defaultValue: 'ws://10.0.2.2:8080/api/v1/chat/ws');
  final String currentUserId;

  ChatRepository({required this.currentUserId});

  void connectWebSocket() {
    _channel = WebSocketChannel.connect(Uri.parse(wsUrl));
  }

  Stream<ChatMessage> get messageStream {
    if (_channel == null) throw Exception('WebSocket not connected');
    return _channel!.stream.map((event) {
      final data = jsonDecode(event);
      return ChatMessage.fromJson(data, currentUserId);
    });
  }

  void sendMessage(ChatMessage message) {
    if (_channel != null) {
      _channel!.sink.add(jsonEncode(message.toJson()));
    } else {
      throw Exception('WebSocket is not connected');
    }
  }

  Future<List<ChatMessage>> getHistory(String matchId) async {
    try {
      final response = await _dio.get('$baseUrl/chat/messages/$matchId');
      if (response.statusCode == 200) {
        final List data = response.data['data'] ?? [];
        return data.map((e) => ChatMessage.fromJson(e, currentUserId)).toList();
      }
      return [];
    } catch (e) {
      throw Exception('Failed to load chat history: $e');
    }
  }

  void disconnect() {
    _channel?.sink.close();
  }
}
