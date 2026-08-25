// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';
import 'package:qlove/core/models/user_model.dart';

abstract class DiscoverState extends Equatable {
  const DiscoverState();

  @override
  List<Object?> get props => [];
}

class DiscoverInitial extends DiscoverState {}

class DiscoverLoading extends DiscoverState {}

class DiscoverLoaded extends DiscoverState {
  final List<UserModel> profiles;
  final bool hasReachedMax;

  const DiscoverLoaded({
    required this.profiles,
    this.hasReachedMax = false,
  });

  DiscoverLoaded copyWith({
    List<UserModel>? profiles,
    bool? hasReachedMax,
  }) {
    return DiscoverLoaded(
      profiles: profiles ?? this.profiles,
      hasReachedMax: hasReachedMax ?? this.hasReachedMax,
    );
  }

  @override
  List<Object?> get props => [profiles, hasReachedMax];
}

class DiscoverMatch extends DiscoverState {
  final UserModel matchedUser;
  final List<UserModel> remainingProfiles;

  const DiscoverMatch({
    required this.matchedUser,
    required this.remainingProfiles,
  });

  @override
  List<Object?> get props => [matchedUser, remainingProfiles];
}

class DiscoverError extends DiscoverState {
  final String message;

  const DiscoverError(this.message);

  @override
  List<Object?> get props => [message];
}
