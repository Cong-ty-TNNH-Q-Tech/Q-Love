import 'package:flutter/material.dart';
import 'package:rive/rive.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

class LoveIslandScreen extends StatefulWidget {
  const LoveIslandScreen({Key? key}) : super(key: key);

  @override
  _LoveIslandScreenState createState() => _LoveIslandScreenState();
}

class _LoveIslandScreenState extends State<LoveIslandScreen> {
  // A placeholder Rive animation URL since we don't have a real island asset yet
  final String _riveUrl = 'https://cdn.rive.app/animations/vehicles.riv';
  
  StateMachineController? _controller;
  
  void _onRiveInit(Artboard artboard) {
    // For a real island asset, we would use its specific state machine name.
    // e.g. StateMachineController.fromArtboard(artboard, 'IslandStateMachine');
    _controller = StateMachineController.fromArtboard(
      artboard,
      'bumpy', // state machine for the placeholder vehicle
    );
    
    if (_controller != null) {
      artboard.addController(_controller!);
    }
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F1A), // Gen-Z Dark-first
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: Text(
          AppLocalizations.of(context)!.loveIslandTitle,
          style: const TextStyle(fontWeight: FontWeight.bold, shadows: [
            Shadow(color: Colors.black54, blurRadius: 4, offset: Offset(1, 1))
          ]),
        ),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white, shadows: [
            Shadow(color: Colors.black54, blurRadius: 4, offset: Offset(1, 1))
          ]),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: Stack(
        fit: StackFit.expand,
        children: [
          // Rive Animation Background
          RiveAnimation.network(
            _riveUrl,
            fit: BoxFit.cover,
            onInit: _onRiveInit,
            alignment: Alignment.center,
          ),
          
          // Gradient Overlay to ensure text readability at the bottom
          Positioned(
            bottom: 0,
            left: 0,
            right: 0,
            height: 200,
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.bottomCenter,
                  end: Alignment.topCenter,
                  colors: [
                    const Color(0xFF0F0F1A),
                    const Color(0xFF0F0F1A).withOpacity(0.0),
                  ],
                ),
              ),
            ),
          ),
          
          // UI Overlays
          Positioned(
            bottom: 40,
            left: 24,
            right: 24,
            child: Column(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: Colors.white.withOpacity(0.2)),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.2),
                        blurRadius: 10,
                      )
                    ]
                  ),
                  child: Text(
                    AppLocalizations.of(context)!.interactWithIsland,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    _buildActionButton(Icons.favorite, Colors.pinkAccent),
                    _buildActionButton(Icons.star, Colors.amber),
                    _buildActionButton(Icons.home, Colors.cyan),
                  ],
                ),
              ],
            ),
          )
        ],
      ),
    );
  }

  Widget _buildActionButton(IconData icon, Color color) {
    return FloatingActionButton(
      heroTag: null, // Avoid hero tag conflicts
      onPressed: () {
        // Handle interaction, maybe trigger a Rive boolean/trigger input
        if (_controller != null) {
          // E.g. _controller!.findInput<bool>('isDay')?.value = false;
        }
      },
      backgroundColor: Colors.white.withOpacity(0.15),
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: BorderSide(color: color.withOpacity(0.5)),
      ),
      child: Icon(icon, color: color, size: 28),
    );
  }
}
