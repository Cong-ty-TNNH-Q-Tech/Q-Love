// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:qlove/features/auth/data/auth_repository.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _authRepository;
  String? _accessToken;

  String? get currentAccessToken => _accessToken;

  AuthBloc({required AuthRepository authRepository})
      : _authRepository = authRepository,
        super(AuthInitial()) {
    on<SendOtpRequested>(_onSendOtpRequested);
    on<VerifyOtpRequested>(_onVerifyOtpRequested);
    on<UpdateProfileRequested>(_onUpdateProfileRequested);
    on<TokenRefreshed>(_onTokenRefreshed);
    on<LogoutRequested>(_onLogoutRequested);
  }

  Future<void> _onSendOtpRequested(
    SendOtpRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      await _authRepository.sendOtp(event.phone);
      emit(AuthOtpSent(event.phone));
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> _onVerifyOtpRequested(
    VerifyOtpRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      final result = await _authRepository.verifyOtp(event.phone, event.otp);
      
      _accessToken = result['access_token'];
      final isNewUser = result['is_new_user'] as bool;
      final user = result['user'];

      if (isNewUser) {
        emit(AuthProfileNeedsCreation(_accessToken!, user));
      } else {
        emit(AuthAuthenticated(_accessToken!, user));
      }
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> _onUpdateProfileRequested(
    UpdateProfileRequested event,
    Emitter<AuthState> emit,
  ) async {
    if (state is! AuthProfileNeedsCreation && state is! AuthAuthenticated) {
      emit(const AuthError('User is not authenticated'));
      return;
    }

    final currentState = state;
    emit(AuthLoading());
    
    try {
      final updatedUser = await _authRepository.updateProfile(
        name: event.name,
        gender: event.gender,
        avatarUrl: event.avatarUrl,
        zodiac: event.zodiac,
      );

      emit(AuthAuthenticated(_accessToken!, updatedUser));
    } catch (e) {
      emit(AuthError(e.toString()));
      // Revert to previous state if needed
      if (currentState is AuthProfileNeedsCreation) {
        emit(currentState);
      } else if (currentState is AuthAuthenticated) {
        emit(currentState);
      }
    }
  }

  void _onTokenRefreshed(TokenRefreshed event, Emitter<AuthState> emit) {
    _accessToken = event.newAccessToken;
    if (state is AuthAuthenticated) {
      emit(AuthAuthenticated(_accessToken!, (state as AuthAuthenticated).user));
    } else if (state is AuthProfileNeedsCreation) {
      emit(AuthProfileNeedsCreation(_accessToken!, (state as AuthProfileNeedsCreation).tempUser));
    }
  }

  Future<void> _onLogoutRequested(LogoutRequested event, Emitter<AuthState> emit) async {
    emit(AuthLoading());
    try {
      await _authRepository.logout();
      _accessToken = null;
      emit(AuthUnauthenticated());
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }
}
