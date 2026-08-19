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
  final int? age;
  final double? distanceKm;
  final int? spiritualCompatibility; // %
  final String? quote;

  const UserModel({
    required this.id,
    this.phone,
    this.name,
    this.dob,
    this.gender,
    this.level,
    this.avatarUrl,
    this.zodiac,
    this.age,
    this.distanceKm,
    this.spiritualCompatibility,
    this.quote,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String? ?? '',
      phone: json['phone'] as String?,
      name: json['name'] as String?,
      dob: json['dob'] as String?,
      gender: json['gender'] as String?,
      level: json['level'] as int?,
      avatarUrl: json['avatar_url'] as String?,
      zodiac: json['zodiac'] as String?,
      age: json['age'] as int?,
      distanceKm: (json['distance_km'] as num?)?.toDouble(),
      spiritualCompatibility: json['spiritual_compatibility'] as int?,
      quote: json['quote'] as String?,
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
      'age': age,
      'distance_km': distanceKm,
      'spiritual_compatibility': spiritualCompatibility,
      'quote': quote,
    };
  }

  @override
  List<Object?> get props => [
        id,
        phone,
        name,
        dob,
        gender,
        level,
        avatarUrl,
        zodiac,
        age,
        distanceKm,
        spiritualCompatibility,
        quote,
      ];
}
