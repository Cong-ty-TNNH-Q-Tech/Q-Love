// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';
import 'package:qlove/core/models/user_model.dart';

abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitial extends AuthState {}

class AuthLoading extends AuthState {}

class AuthOtpSent extends AuthState {
  final String phone;

  const AuthOtpSent(this.phone);

  @override
  List<Object?> get props => [phone];
}

class AuthProfileNeedsCreation extends AuthState {
  final String accessToken;
  final UserModel tempUser;

  const AuthProfileNeedsCreation(this.accessToken, this.tempUser);

  @override
  List<Object?> get props => [accessToken, tempUser];
}

class AuthAuthenticated extends AuthState {
  final String accessToken;
  final UserModel user;

  const AuthAuthenticated(this.accessToken, this.user);

  @override
  List<Object?> get props => [accessToken, user];
}

class AuthUnauthenticated extends AuthState {}

class AuthError extends AuthState {
  final String message;

  const AuthError(this.message);

  @override
  List<Object?> get props => [message];
}
