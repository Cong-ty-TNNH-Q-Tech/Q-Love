// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

// --- Events ---
abstract class WingmanEvent extends Equatable {
  const WingmanEvent();

  @override
  List<Object?> get props => [];
}

class LoadWingmanDashboard extends WingmanEvent {}

class MatchFriendEvent extends WingmanEvent {
  final String targetId;
  final String friendId;

  const MatchFriendEvent({required this.targetId, required this.friendId});

  @override
  List<Object?> get props => [targetId, friendId];
}

// --- States ---
abstract class WingmanState extends Equatable {
  const WingmanState();

  @override
  List<Object?> get props => [];
}

class WingmanInitial extends WingmanState {}

class WingmanLoading extends WingmanState {}

class WingmanDashboardLoaded extends WingmanState {
  final int totalMatches;
  final double successRate;
  final double totalCommission;
  final List<double> chartData;

  const WingmanDashboardLoaded({
    required this.totalMatches,
    required this.successRate,
    required this.totalCommission,
    required this.chartData,
  });

  @override
  List<Object?> get props => [totalMatches, successRate, totalCommission, chartData];
}

class WingmanError extends WingmanState {
  final String message;

  const WingmanError(this.message);

  @override
  List<Object?> get props => [message];
}

class MatchmakerSuccess extends WingmanState {}

// --- Bloc ---
class WingmanBloc extends Bloc<WingmanEvent, WingmanState> {
  WingmanBloc() : super(WingmanInitial()) {
    on<LoadWingmanDashboard>(_onLoadDashboard);
    on<MatchFriendEvent>(_onMatchFriend);
  }

  Future<void> _onLoadDashboard(LoadWingmanDashboard event, Emitter<WingmanState> emit) async {
    emit(WingmanLoading());
    try {
      // TODO: Replace with real API call
      await Future.delayed(const Duration(seconds: 1));
      
      emit(const WingmanDashboardLoaded(
        totalMatches: 42,
        successRate: 0.85,
        totalCommission: 12500, // Xu
        chartData: [1000, 1500, 2000, 800, 3000, 2500, 1700], // Last 7 days
      ));
    } catch (e) {
      emit(WingmanError(e.toString()));
    }
  }

  Future<void> _onMatchFriend(MatchFriendEvent event, Emitter<WingmanState> emit) async {
    emit(WingmanLoading());
    try {
      // TODO: Replace with real API call
      await Future.delayed(const Duration(seconds: 1));
      emit(MatchmakerSuccess());
    } catch (e) {
      emit(WingmanError(e.toString()));
    }
  }
}
