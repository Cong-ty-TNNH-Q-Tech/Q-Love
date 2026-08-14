// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/vibe_check_cubit.dart';
import '../bloc/vibe_check_state.dart';
import 'widgets/spotify_card_widget.dart';

class VibeCheckScreen extends StatefulWidget {
  const VibeCheckScreen({super.key});

  @override
  State<VibeCheckScreen> createState() => _VibeCheckScreenState();
}

class _VibeCheckScreenState extends State<VibeCheckScreen> {
  @override
  Widget build(BuildContext context) {
    // Vibe Check UI enforces a deep dark mode
    return Theme(
      data: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: const Color(0xFF0F0F1A),
      ),
      child: Scaffold(
        appBar: AppBar(
          backgroundColor: Colors.transparent,
          elevation: 0,
          title: const Text(
            'Midnight Vibe Check',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: Colors.white,
              shadows: [
                Shadow(
                  color: Colors.purpleAccent,
                  blurRadius: 10,
                )
              ]
            ),
          ),
          centerTitle: true,
        ),
        body: BlocProvider(
          create: (context) => VibeCheckCubit()..fetchVibeTracks(),
          child: Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text(
                    "Discover who's listening with you...",
                    style: TextStyle(color: Colors.white70, fontSize: 16),
                  ),
                  const SizedBox(height: 32),
                  BlocBuilder<VibeCheckCubit, VibeCheckState>(
                    builder: (context, state) {
                      if (state is VibeCheckLoading || state is VibeCheckInitial) {
                        return const Center(child: CircularProgressIndicator(color: Colors.purpleAccent));
                      } else if (state is VibeCheckError) {
                        return Text(
                          state.message,
                          style: const TextStyle(color: Colors.redAccent, fontSize: 16),
                          textAlign: TextAlign.center,
                        );
                      } else if (state is VibeCheckLoaded) {
                        final tracks = state.tracks;
                        if (tracks.isNotEmpty) {
                          return Dismissible(
                            key: Key(tracks[0]['title']!),
                            onDismissed: (direction) {
                              context.read<VibeCheckCubit>().removeTrack(0);
                            },
                            child: SpotifyCardWidget(
                              title: tracks[0]['title']!,
                              artist: tracks[0]['artist']!,
                              coverUrl: tracks[0]['coverUrl']!,
                              previewUrl: tracks[0]['previewUrl']!,
                            ),
                          );
                        } else {
                          return const Text(
                            "No more vibes tonight. Check back tomorrow!",
                            style: TextStyle(color: Colors.white54, fontSize: 18),
                            textAlign: TextAlign.center,
                          );
                        }
                      }
                      return const SizedBox.shrink();
                    },
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
