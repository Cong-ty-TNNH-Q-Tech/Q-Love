// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

part of 'shame_wall_bloc.dart';

abstract class ShameWallState extends Equatable {
  const ShameWallState();
  
  @override
  List<Object> get props => [];
}

class ShameWallInitial extends ShameWallState {}

class ShameWallLoading extends ShameWallState {}

class ShameWallLoaded extends ShameWallState {
  final List<ShameModel> shames;

  const ShameWallLoaded({required this.shames});

  @override
  List<Object> get props => [shames];
}

class ShameWallError extends ShameWallState {
  final String message;

  const ShameWallError({required this.message});

  @override
  List<Object> get props => [message];
}
