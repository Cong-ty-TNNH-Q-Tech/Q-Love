// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:qlove/core/models/user_model.dart';
import 'package:qlove/features/discover/bloc/discover_event.dart';
import 'package:qlove/features/discover/bloc/discover_state.dart';
import 'package:qlove/features/discover/data/discover_repository.dart';

class DiscoverBloc extends Bloc<DiscoverEvent, DiscoverState> {
  final DiscoverRepository discoverRepository;

  DiscoverBloc({required this.discoverRepository}) : super(DiscoverInitial()) {
    on<FetchFeedRequested>(_onFetchFeedRequested);
    on<SwipeActionRequested>(_onSwipeActionRequested);
    on<ResumeDiscoverRequested>(_onResumeDiscoverRequested);
  }

  Future<void> _onFetchFeedRequested(
      FetchFeedRequested event, Emitter<DiscoverState> emit) async {
    try {
      if (state is DiscoverInitial || event.isRefresh) {
        emit(DiscoverLoading());
      }
      
      final feed = await discoverRepository.getFeed(filter: event.filter);
      
      if (feed.isEmpty) {
        emit(const DiscoverLoaded(profiles: [], hasReachedMax: true));
      } else {
        emit(DiscoverLoaded(profiles: feed, hasReachedMax: false));
      }
    } catch (e) {
      emit(DiscoverError(e.toString()));
    }
  }

  Future<void> _onSwipeActionRequested(
      SwipeActionRequested event, Emitter<DiscoverState> emit) async {
    final currentState = state;
    if (currentState is DiscoverLoaded) {
      try {
        final remainingProfiles = List<UserModel>.from(currentState.profiles)
          ..removeWhere((p) => p.id == event.targetId);
        
        // Optimistic update
        emit(DiscoverLoaded(
            profiles: remainingProfiles,
            hasReachedMax: remainingProfiles.isEmpty));

        final isMatch = await discoverRepository.swipe(event.targetId, event.action);
        
        if (isMatch && event.action == 'like') {
          // If match, emit DiscoverMatch
          final matchedUser = currentState.profiles.firstWhere((p) => p.id == event.targetId);
          emit(DiscoverMatch(
            matchedUser: matchedUser,
            remainingProfiles: remainingProfiles,
          ));
        }
      } catch (e) {
        // Handle error if necessary, maybe revert optimistic update
        emit(DiscoverError(e.toString()));
      }
    }
  }

  void _onResumeDiscoverRequested(
      ResumeDiscoverRequested event, Emitter<DiscoverState> emit) {
    if (state is DiscoverMatch) {
      final matchState = state as DiscoverMatch;
      emit(DiscoverLoaded(
          profiles: matchState.remainingProfiles,
          hasReachedMax: matchState.remainingProfiles.isEmpty));
    }
  }
}
