// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:async';
import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/minigame_api_client.dart';

// Events
abstract class CardStealEvent extends Equatable {
  const CardStealEvent();
  @override
  List<Object> get props => [];
}

class StartStealGame extends CardStealEvent {
  final String defenderId;
  final String targetCardId;
  const StartStealGame(this.defenderId, this.targetCardId);
}

class TapCard extends CardStealEvent {}

class TimerTick extends CardStealEvent {}

class SubmitResult extends CardStealEvent {}

// States
abstract class CardStealState extends Equatable {
  const CardStealState();
  @override
  List<Object?> get props => [];
}

class CardStealInitial extends CardStealState {}

class CardStealLoading extends CardStealState {}

class CardStealPlaying extends CardStealState {
  final int taps;
  final int timeLeft;
  final String stealId;
  
  const CardStealPlaying({required this.taps, required this.timeLeft, required this.stealId});
  
  @override
  List<Object> get props => [taps, timeLeft, stealId];
}

class CardStealWon extends CardStealState {}
class CardStealLost extends CardStealState {}
class CardStealError extends CardStealState {
  final String message;
  const CardStealError(this.message);
}

// BLoC
class CardStealBloc extends Bloc<CardStealEvent, CardStealState> {
  final MinigameApiClient apiClient;
  Timer? _timer;
  
  static const int requiredTaps = 50;
  static const int gameDurationSeconds = 10;
  
  CardStealBloc({required this.apiClient}) : super(CardStealInitial()) {
    on<StartStealGame>(_onStartStealGame);
    on<TapCard>(_onTapCard);
    on<TimerTick>(_onTimerTick);
    on<SubmitResult>(_onSubmitResult);
  }
  
  Future<void> _onStartStealGame(StartStealGame event, Emitter<CardStealState> emit) async {
    emit(CardStealLoading());
    try {
      final stealId = await apiClient.initSteal(event.defenderId, event.targetCardId);
      emit(CardStealPlaying(taps: 0, timeLeft: gameDurationSeconds, stealId: stealId));
      
      _timer?.cancel();
      _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
        add(TimerTick());
      });
    } catch (e) {
      emit(CardStealError(e.toString()));
    }
  }
  
  void _onTapCard(TapCard event, Emitter<CardStealState> emit) {
    if (state is CardStealPlaying) {
      final currentState = state as CardStealPlaying;
      final newTaps = currentState.taps + 1;
      
      if (newTaps >= requiredTaps) {
        _timer?.cancel();
        emit(CardStealPlaying(taps: newTaps, timeLeft: currentState.timeLeft, stealId: currentState.stealId));
        add(SubmitResult());
      } else {
        emit(CardStealPlaying(taps: newTaps, timeLeft: currentState.timeLeft, stealId: currentState.stealId));
      }
    }
  }
  
  void _onTimerTick(TimerTick event, Emitter<CardStealState> emit) {
    if (state is CardStealPlaying) {
      final currentState = state as CardStealPlaying;
      if (currentState.timeLeft > 0) {
        emit(CardStealPlaying(taps: currentState.taps, timeLeft: currentState.timeLeft - 1, stealId: currentState.stealId));
      } else {
        _timer?.cancel();
        add(SubmitResult());
      }
    }
  }
  
  Future<void> _onSubmitResult(SubmitResult event, Emitter<CardStealState> emit) async {
    if (state is CardStealPlaying) {
      final currentState = state as CardStealPlaying;
      final isWin = currentState.taps >= requiredTaps;
      
      emit(CardStealLoading());
      try {
        await apiClient.submitStealResult(currentState.stealId, isWin);
        if (isWin) {
          emit(CardStealWon());
        } else {
          emit(CardStealLost());
        }
      } catch (e) {
        emit(CardStealError(e.toString()));
      }
    }
  }
  
  @override
  Future<void> close() {
    _timer?.cancel();
    return super.close();
  }
}
