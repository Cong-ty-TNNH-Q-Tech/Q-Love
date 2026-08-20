// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

class ShameModel extends Equatable {
  final String id;
  final String userName;
  final String reason;
  final int tomatoes;
  final DateTime expiresAt;

  const ShameModel({
    required this.id,
    required this.userName,
    required this.reason,
    required this.tomatoes,
    required this.expiresAt,
  });

  factory ShameModel.fromJson(Map<String, dynamic> json) {
    return ShameModel(
      id: json['id'] as String,
      userName: json['user_name'] as String? ?? 'Ẩn danh',
      reason: json['reason'] as String? ?? 'Chưa rõ',
      tomatoes: json['tomato_count'] as int? ?? 0,
      expiresAt: json['expires_at'] != null 
          ? DateTime.parse(json['expires_at'])
          : DateTime.now().add(const Duration(hours: 24)),
    );
  }

  ShameModel copyWith({
    String? id,
    String? userName,
    String? reason,
    int? tomatoes,
    DateTime? expiresAt,
  }) {
    return ShameModel(
      id: id ?? this.id,
      userName: userName ?? this.userName,
      reason: reason ?? this.reason,
      tomatoes: tomatoes ?? this.tomatoes,
      expiresAt: expiresAt ?? this.expiresAt,
    );
  }

  @override
  List<Object> get props => [id, userName, reason, tomatoes, expiresAt];
}
