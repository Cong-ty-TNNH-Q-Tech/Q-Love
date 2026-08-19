// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class SendOtpRequested extends AuthEvent {
  final String phone;

  const SendOtpRequested(this.phone);

  @override
  List<Object?> get props => [phone];
}

class VerifyOtpRequested extends AuthEvent {
  final String phone;
  final String otp;

  const VerifyOtpRequested({required this.phone, required this.otp});

  @override
  List<Object?> get props => [phone, otp];
}

class UpdateProfileRequested extends AuthEvent {
  final String name;
  final String gender;
  final String? avatarUrl;
  final String? zodiac;

  const UpdateProfileRequested({
    required this.name,
    required this.gender,
    this.avatarUrl,
    this.zodiac,
  });

  @override
  List<Object?> get props => [name, gender, avatarUrl, zodiac];
}

class TokenRefreshed extends AuthEvent {
  final String newAccessToken;

  const TokenRefreshed(this.newAccessToken);

  @override
  List<Object?> get props => [newAccessToken];
}

class LogoutRequested extends AuthEvent {}
