// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/features/auth/bloc/auth_bloc.dart';
import 'package:qlove/features/auth/bloc/auth_event.dart';
import 'package:qlove/features/auth/bloc/auth_state.dart';
import 'package:qlove/features/auth/data/auth_repository.dart';

class MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late AuthBloc authBloc;
  late MockAuthRepository mockAuthRepository;

  setUp(() {
    mockAuthRepository = MockAuthRepository();
    authBloc = AuthBloc(authRepository: mockAuthRepository);
  });

  tearDown(() {
    authBloc.close();
  });

  group('AuthBloc', () {
    test('initial state is AuthInitial', () {
      expect(authBloc.state, isA<AuthInitial>());
    });

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthOtpSent] when SendOtpRequested succeeds',
      build: () {
        when(() => mockAuthRepository.sendOtp('0901234567'))
            .thenAnswer((_) async {});
        return authBloc;
      },
      act: (bloc) => bloc.add(const SendOtpRequested('0901234567')),
      expect: () => [
        isA<AuthLoading>(),
        isA<AuthOtpSent>().having((state) => state.phone, 'phone', '0901234567'),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthAuthenticated] when VerifyOtpRequested succeeds and user is not new',
      build: () {
        final mockUser = const UserModel(id: '1', phone: '0901234567');
        when(() => mockAuthRepository.verifyOtp('0901234567', '123456'))
            .thenAnswer((_) async => {
                  'access_token': 'token123',
                  'is_new_user': false,
                  'user': mockUser,
                });
        return authBloc;
      },
      act: (bloc) => bloc.add(const VerifyOtpRequested(phone: '0901234567', otp: '123456')),
      expect: () => [
        isA<AuthLoading>(),
        isA<AuthAuthenticated>()
            .having((state) => state.accessToken, 'token', 'token123')
            .having((state) => state.user.phone, 'user.phone', '0901234567'),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthProfileNeedsCreation] when VerifyOtpRequested succeeds and user is new',
      build: () {
        final mockUser = const UserModel(id: '1', phone: '0901234567');
        when(() => mockAuthRepository.verifyOtp('0901234567', '123456'))
            .thenAnswer((_) async => {
                  'access_token': 'token123',
                  'is_new_user': true,
                  'user': mockUser,
                });
        return authBloc;
      },
      act: (bloc) => bloc.add(const VerifyOtpRequested(phone: '0901234567', otp: '123456')),
      expect: () => [
        isA<AuthLoading>(),
        isA<AuthProfileNeedsCreation>()
            .having((state) => state.accessToken, 'token', 'token123')
            .having((state) => state.tempUser.phone, 'user.phone', '0901234567'),
      ],
    );
  });
}
