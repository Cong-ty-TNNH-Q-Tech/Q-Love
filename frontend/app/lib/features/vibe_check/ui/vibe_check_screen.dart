import 'package:flutter/material.dart';
import 'widgets/spotify_card_widget.dart';

class VibeCheckScreen extends StatefulWidget {
  const VibeCheckScreen({super.key});

  @override
  State<VibeCheckScreen> createState() => _VibeCheckScreenState();
}

class _VibeCheckScreenState extends State<VibeCheckScreen> {
  // Mock data representing Spotify Tracks fetched from backend API
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
        body: Center(
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
                if (tracks.isNotEmpty)
                  Dismissible(
                    key: Key(tracks[0]['title']!),
                    onDismissed: (direction) {
                      setState(() {
                        tracks.removeAt(0);
                      });
                    },
                    child: SpotifyCardWidget(
                      title: tracks[0]['title']!,
                      artist: tracks[0]['artist']!,
                      coverUrl: tracks[0]['coverUrl']!,
                      previewUrl: tracks[0]['previewUrl']!,
                    ),
                  )
                else
                  const Text(
                    "No more vibes tonight. Check back tomorrow!",
                    style: TextStyle(color: Colors.white54, fontSize: 18),
                    textAlign: TextAlign.center,
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
