// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

part 'shame_wall_event.dart';
part 'shame_wall_state.dart';

class ShameWallBloc extends Bloc<ShameWallEvent, ShameWallState> {
  ShameWallBloc() : super(ShameWallInitial()) {
    on<LoadShameWall>(_onLoadShameWall);
    on<ThrowTomato>(_onThrowTomato);
  }

  void _onLoadShameWall(LoadShameWall event, Emitter<ShameWallState> emit) async {
    emit(ShameWallLoading());
    try {
      // Simulate API call
      await Future.delayed(const Duration(seconds: 1));
      final shames = [
        ShameModel(id: '1', userName: 'Kẻ Bội Tín', reason: 'Bắt cá 2 tay', tomatoes: 120, expiresAt: DateTime.now().add(const Duration(hours: 12))),
        ShameModel(id: '2', userName: 'Tra Nam', reason: 'Ghosting sau chốt kèo', tomatoes: 45, expiresAt: DateTime.now().add(const Duration(hours: 5))),
      ];
      emit(ShameWallLoaded(shames: shames));
    } catch (e) {
      emit(ShameWallError(message: e.toString()));
    }
  }

  void _onThrowTomato(ThrowTomato event, Emitter<ShameWallState> emit) async {
    if (state is ShameWallLoaded) {
      final currentState = state as ShameWallLoaded;
      try {
        // Simulate API call to throw tomato
        await Future.delayed(const Duration(milliseconds: 300));
        
        // Optimistic update
        final updatedShames = currentState.shames.map((s) {
          if (s.id == event.shameId) {
            return s.copyWith(tomatoes: s.tomatoes + 1);
          }
          return s;
        }).toList();

        emit(ShameWallLoaded(shames: updatedShames));
      } catch (e) {
        // Fallback to current state or show error
        emit(ShameWallError(message: 'Không đủ Xu để ném cà chua'));
      }
    }
  }
}

class ShameModel extends Equatable {
  final String id;
  final String userName;
  final String reason;
  final int tomatoes;
  final DateTime expiresAt;

  const ShameModel({
    required this.id,
    required this.userName,
    required this.reason,
    required this.tomatoes,
    required this.expiresAt,
  });

  ShameModel copyWith({
    String? id,
    String? userName,
    String? reason,
    int? tomatoes,
    DateTime? expiresAt,
  }) {
    return ShameModel(
      id: id ?? this.id,
      userName: userName ?? this.userName,
      reason: reason ?? this.reason,
      tomatoes: tomatoes ?? this.tomatoes,
      expiresAt: expiresAt ?? this.expiresAt,
    );
  }

  @override
  List<Object> get props => [id, userName, reason, tomatoes, expiresAt];
}
