// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

class CourtCaseModel extends Equatable {
  final String id;
  final String defendantNameMasked;
  final String reason;
  final String status;
  final int voteCount;
  final DateTime createdAt;

  const CourtCaseModel({
    required this.id,
    required this.defendantNameMasked,
    required this.reason,
    required this.status,
    required this.voteCount,
    required this.createdAt,
  });

  factory CourtCaseModel.fromJson(Map<String, dynamic> json) {
    return CourtCaseModel(
      id: json['id'] as String,
      defendantNameMasked: json['defendant_name_masked'] as String? ?? 'Ẩn danh ***',
      reason: json['reason'] as String? ?? 'Không rõ lý do',
      status: json['status'] as String? ?? 'voting',
      voteCount: json['vote_count'] as int? ?? 0,
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
    );
  }

  @override
  List<Object?> get props => [id, defendantNameMasked, reason, status, voteCount, createdAt];
}
