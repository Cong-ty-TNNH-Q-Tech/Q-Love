// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'vibe_check_state.dart';
import '../data/vibe_repository.dart';

class VibeCheckCubit extends Cubit<VibeCheckState> {
  final VibeRepository repository;

  VibeCheckCubit({required this.repository}) : super(VibeCheckInitial());

  Future<void> fetchVibeTracks() async {
    emit(VibeCheckLoading());
    try {
      final track = await repository.getCurrentVibeTrack();
      final List<Map<String, String>> tracks = [track];
      emit(VibeCheckLoaded(tracks));
    } catch (e) {
      emit(VibeCheckError("Failed to load vibe tracks: $e"));
    }
  }

  void removeTrack(int index) {
    if (state is VibeCheckLoaded) {
      final currentState = state as VibeCheckLoaded;
      final newTracks = List<Map<String, String>>.from(currentState.tracks);
      if (index >= 0 && index < newTracks.length) {
        newTracks.removeAt(index);
        emit(VibeCheckLoaded(newTracks));
      }
    }
  }

  Future<void> matchTrack(int index) async {
    if (state is VibeCheckLoaded) {
      final currentState = state as VibeCheckLoaded;
      if (index >= 0 && index < currentState.tracks.length) {
        final trackId = currentState.tracks[index]['id'];
        if (trackId != null && trackId.isNotEmpty) {
          try {
            await repository.matchVibe(trackId);
          } catch (e) {
            // Ignore for now, or emit some error state
          }
        }
        removeTrack(index);
      }
    }
  }
}
