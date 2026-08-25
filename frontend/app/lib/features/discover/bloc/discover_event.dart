// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:equatable/equatable.dart';

abstract class DiscoverEvent extends Equatable {
  const DiscoverEvent();

  @override
  List<Object?> get props => [];
}

class FetchFeedRequested extends DiscoverEvent {
  final String filter;
  final bool isRefresh;

  const FetchFeedRequested({this.filter = 'default', this.isRefresh = false});

  @override
  List<Object?> get props => [filter, isRefresh];
}

class SwipeActionRequested extends DiscoverEvent {
  final String targetId;
  final String action; // 'like' or 'pass'

  const SwipeActionRequested({required this.targetId, required this.action});

  @override
  List<Object?> get props => [targetId, action];
}

class ResumeDiscoverRequested extends DiscoverEvent {}
