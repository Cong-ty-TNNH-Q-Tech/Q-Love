// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

class ChatMessage {
  final String id;
  final String matchId;
  final String senderId;
  final String targetId;
  final String type; // 'text', 'locket'
  final String content;
  final DateTime createdAt;
  final bool isMine;

  ChatMessage({
    required this.id,
    required this.matchId,
    required this.senderId,
    required this.targetId,
    required this.type,
    required this.content,
    required this.createdAt,
    required this.isMine,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json, String currentUserId) {
    return ChatMessage(
      id: json['id'] ?? '',
      matchId: json['match_id'] ?? '',
      senderId: json['sender_id'] ?? '',
      targetId: json['target_id'] ?? '',
      type: json['type'] ?? 'text',
      content: json['content'] ?? '',
      createdAt: json['created_at'] != null ? DateTime.parse(json['created_at']) : DateTime.now(),
      isMine: json['sender_id'] == currentUserId,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'match_id': matchId,
      'target_id': targetId,
      'type': type,
      'content': content,
    };
  }
}
