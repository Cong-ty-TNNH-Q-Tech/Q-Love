// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:qlove/core/models/court_case_model.dart';
import 'package:qlove/features/court/bloc/court_event.dart';
import 'package:qlove/features/court/bloc/court_state.dart';
import 'package:qlove/features/court/data/court_repository.dart';

class CourtBloc extends Bloc<CourtEvent, CourtState> {
  final CourtRepository courtRepository;

  CourtBloc({required this.courtRepository}) : super(CourtInitial()) {
    on<FetchCasesRequested>(_onFetchCasesRequested);
    on<VoteActionRequested>(_onVoteActionRequested);
  }

  Future<void> _onFetchCasesRequested(
      FetchCasesRequested event, Emitter<CourtState> emit) async {
    try {
      if (state is CourtInitial || event.isRefresh) {
        emit(CourtLoading());
      }
      
      final cases = await courtRepository.getCases();
      
      if (cases.isEmpty) {
        emit(const CourtLoaded(cases: [], hasReachedMax: true));
      } else {
        emit(CourtLoaded(cases: cases, hasReachedMax: false));
      }
    } catch (e) {
      emit(CourtError(e.toString()));
    }
  }

  Future<void> _onVoteActionRequested(
      VoteActionRequested event, Emitter<CourtState> emit) async {
    final currentState = state;
    if (currentState is CourtLoaded) {
      try {
        final remainingCases = List<CourtCaseModel>.from(currentState.cases)
          ..removeWhere((c) => c.id == event.caseId);
        
        // Optimistic update
        emit(CourtLoaded(
            cases: remainingCases,
            hasReachedMax: remainingCases.isEmpty));

        await courtRepository.vote(event.caseId, event.voteType);
        
      } catch (e) {
        // Handle error, maybe revert optimistic update
        emit(CourtError(e.toString()));
      }
    }
  }
}
