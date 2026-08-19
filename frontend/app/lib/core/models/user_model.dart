// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

class UserModel extends Equatable {
  final String id;
  final String? phone;
  final String? name;
  final String? dob;
  final String? gender;
  final int? level;
  final String? avatarUrl;
  final String? zodiac;

  const UserModel({
    required this.id,
    this.phone,
    this.name,
    this.dob,
    this.gender,
    this.level,
    this.avatarUrl,
    this.zodiac,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] ?? '',
      phone: json['phone'],
      name: json['name'],
      dob: json['dob'],
      gender: json['gender'],
      level: json['level'],
      avatarUrl: json['avatar_url'],
      zodiac: json['zodiac'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'phone': phone,
      'name': name,
      'dob': dob,
      'gender': gender,
      'level': level,
      'avatar_url': avatarUrl,
      'zodiac': zodiac,
    };
  }

  @override
  List<Object?> get props => [id, phone, name, dob, gender, level, avatarUrl, zodiac];
}
