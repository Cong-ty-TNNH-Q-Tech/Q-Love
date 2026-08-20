// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:qlove/features/drama/models/shame_model.dart';
import 'package:qlove/features/drama/repositories/shame_repository.dart';

part 'shame_wall_event.dart';
part 'shame_wall_state.dart';

class ShameWallBloc extends Bloc<ShameWallEvent, ShameWallState> {
  final ShameRepository _repository;

  ShameWallBloc({required ShameRepository repository}) 
      : _repository = repository,
        super(ShameWallInitial()) {
    on<LoadShameWall>(_onLoadShameWall);
    on<ThrowTomato>(_onThrowTomato);
  }

  void _onLoadShameWall(LoadShameWall event, Emitter<ShameWallState> emit) async {
    emit(ShameWallLoading());
    try {
      final shames = await _repository.getActiveShames();
      emit(ShameWallLoaded(shames: shames));
    } catch (e) {
      emit(ShameWallError(message: e.toString()));
    }
  }

  void _onThrowTomato(ThrowTomato event, Emitter<ShameWallState> emit) async {
    if (state is ShameWallLoaded) {
      final currentState = state as ShameWallLoaded;
      
      // Optimistic update
      final updatedShames = currentState.shames.map((s) {
        if (s.id == event.shameId) {
          return s.copyWith(tomatoes: s.tomatoes + 1);
        }
        return s;
      }).toList();
      emit(ShameWallLoaded(shames: updatedShames));

      try {
        await _repository.throwTomato(event.shameId);
      } catch (e) {
        // Rollback optimistic update
        emit(ShameWallLoaded(shames: currentState.shames));
        emit(ShameWallError(message: 'Không đủ Xu hoặc có lỗi xảy ra khi ném cà chua'));
      }
    }
  }
}
