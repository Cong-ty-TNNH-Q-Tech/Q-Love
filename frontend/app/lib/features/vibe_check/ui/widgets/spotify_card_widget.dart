// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'package:flutter/material.dart';
import 'package:just_audio/just_audio.dart';

class SpotifyCardWidget extends StatefulWidget {
  final String title;
  final String artist;
  final String coverUrl;
  final String previewUrl;

  const SpotifyCardWidget({
    super.key,
    required this.title,
    required this.artist,
    required this.coverUrl,
    required this.previewUrl,
  });

  @override
  State<SpotifyCardWidget> createState() => _SpotifyCardWidgetState();
}

class _SpotifyCardWidgetState extends State<SpotifyCardWidget> {
  late final AudioPlayer _player;
  bool _isPlaying = false;

  @override
  void initState() {
    super.initState();
    _player = AudioPlayer();
    _initAudio();
  }

  Future<void> _initAudio() async {
    try {
      if (widget.previewUrl.isNotEmpty) {
        final uri = Uri.tryParse(widget.previewUrl);
        if (uri != null && uri.hasAbsolutePath) {
          await _player.setUrl(widget.previewUrl);
        }
      }
    } catch (e) {
      debugPrint("Error loading audio: $e");
    }
  }

  @override
  void dispose() {
    _player.dispose();
    super.dispose();
  }

  void _togglePlay() {
    setState(() {
      _isPlaying = !_isPlaying;
    });
    if (_isPlaying) {
      _player.play();
    } else {
      _player.pause();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: 400,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        boxShadow: [
          BoxShadow(
            color: Colors.purpleAccent.withOpacity(0.3),
            blurRadius: 20,
            spreadRadius: 2,
            offset: const Offset(0, 10),
          )
        ],
        image: DecorationImage(
          image: NetworkImage(widget.coverUrl),
          fit: BoxFit.cover,
        ),
      ),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(24),
          gradient: LinearGradient(
            colors: [
              Colors.black.withOpacity(0.1),
              Colors.black.withOpacity(0.9),
            ],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              widget.title,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),
            Text(
              widget.artist,
              style: const TextStyle(
                color: Colors.white70,
                fontSize: 16,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                IconButton(
                  onPressed: _togglePlay,
                  icon: Icon(
                    _isPlaying ? Icons.pause_circle_filled : Icons.play_circle_fill,
                    size: 48,
                    color: Colors.greenAccent,
                  ),
                ),
                const SizedBox(width: 16),
                const Expanded(
                  child: LinearProgressIndicator(
                    value: 0.3, // Mock progress
                    color: Colors.greenAccent,
                    backgroundColor: Colors.white24,
                  ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}
