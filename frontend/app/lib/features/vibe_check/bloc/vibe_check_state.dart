// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

abstract class VibeCheckState extends Equatable {
  const VibeCheckState();

  @override
  List<Object> get props => [];
}

class VibeCheckInitial extends VibeCheckState {}

class VibeCheckLoading extends VibeCheckState {}

class VibeCheckLoaded extends VibeCheckState {
  final List<Map<String, String>> tracks;

  const VibeCheckLoaded(this.tracks);

  @override
  List<Object> get props => [tracks];
}

class VibeCheckError extends VibeCheckState {
  final String message;

  const VibeCheckError(this.message);

  @override
  List<Object> get props => [message];
}
