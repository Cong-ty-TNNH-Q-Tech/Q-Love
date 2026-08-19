// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';
import 'package:qlove/core/models/court_case_model.dart';

abstract class CourtState extends Equatable {
  const CourtState();

  @override
  List<Object?> get props => [];
}

class CourtInitial extends CourtState {}

class CourtLoading extends CourtState {}

class CourtLoaded extends CourtState {
  final List<CourtCaseModel> cases;
  final bool hasReachedMax;

  const CourtLoaded({
    required this.cases,
    this.hasReachedMax = false,
  });

  CourtLoaded copyWith({
    List<CourtCaseModel>? cases,
    bool? hasReachedMax,
  }) {
    return CourtLoaded(
      cases: cases ?? this.cases,
      hasReachedMax: hasReachedMax ?? this.hasReachedMax,
    );
  }

  @override
  List<Object?> get props => [cases, hasReachedMax];
}

class CourtError extends CourtState {
  final String message;

  const CourtError(this.message);

  @override
  List<Object?> get props => [message];
}
