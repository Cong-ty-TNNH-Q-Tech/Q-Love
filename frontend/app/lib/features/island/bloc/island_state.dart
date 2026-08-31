import 'package:equatable/equatable.dart';

abstract class IslandState extends Equatable {
  const IslandState();

  @override
  List<Object> get props => [];
}

class IslandInitial extends IslandState {}

class IslandLoaded extends IslandState {}
