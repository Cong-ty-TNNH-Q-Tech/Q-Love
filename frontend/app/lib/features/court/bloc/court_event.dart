// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

abstract class CourtEvent extends Equatable {
  const CourtEvent();

  @override
  List<Object?> get props => [];
}

class FetchCasesRequested extends CourtEvent {
  final bool isRefresh;

  const FetchCasesRequested({this.isRefresh = false});

  @override
  List<Object?> get props => [isRefresh];
}

class VoteActionRequested extends CourtEvent {
  final String caseId;
  final String voteType; // 'guilty' or 'not_guilty'

  const VoteActionRequested({required this.caseId, required this.voteType});

  @override
  List<Object?> get props => [caseId, voteType];
}
