// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

part of 'shame_wall_bloc.dart';

abstract class ShameWallEvent extends Equatable {
  const ShameWallEvent();

  @override
  List<Object> get props => [];
}

class LoadShameWall extends ShameWallEvent {}

class ThrowTomato extends ShameWallEvent {
  final String shameId;

  const ThrowTomato({required this.shameId});

  @override
  List<Object> get props => [shameId];
}
