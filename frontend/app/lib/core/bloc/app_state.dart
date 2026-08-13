// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

part of 'app_bloc.dart';

class AppState extends Equatable {
  final bool isInitialized;

  const AppState({
    this.isInitialized = false,
  });

  AppState copyWith({
    bool? isInitialized,
  }) {
    return AppState(
      isInitialized: isInitialized ?? this.isInitialized,
    );
  }

  @override
  List<Object> get props => [isInitialized];
}
