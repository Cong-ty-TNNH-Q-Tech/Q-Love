// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter_bloc/flutter_bloc.dart';
import 'vibe_check_state.dart';

class VibeCheckCubit extends Cubit<VibeCheckState> {
  VibeCheckCubit() : super(VibeCheckInitial());

  Future<void> fetchVibeTracks() async {
    emit(VibeCheckLoading());
    try {
      // Simulate API call to backend
      await Future.delayed(const Duration(seconds: 1));
      
      final List<Map<String, String>> tracks = [
        {
          "title": "Shape of You",
          "artist": "Ed Sheeran",
          "coverUrl": "https://upload.wikimedia.org/wikipedia/en/b/b4/Shape_Of_You_%28Official_Single_Cover%29_by_Ed_Sheeran.png",
          "previewUrl": "",
        },
        {
          "title": "Blinding Lights",
          "artist": "The Weeknd",
          "coverUrl": "https://upload.wikimedia.org/wikipedia/en/e/e6/The_Weeknd_-_Blinding_Lights.png",
          "previewUrl": "",
        }
      ];
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
}
