import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:rive/rive.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../bloc/island_cubit.dart';

class LoveIslandScreen extends StatelessWidget {
  const LoveIslandScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => IslandCubit(),
      child: const _LoveIslandBody(),
    );
  }
}

class _LoveIslandBody extends StatelessWidget {
  const _LoveIslandBody({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final cubit = context.read<IslandCubit>();
    
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
          RiveAnimation.asset(
            'assets/animations/love_island.riv',
            fit: BoxFit.cover,
            onInit: cubit.onRiveInit,
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
                    _buildActionButton(Icons.favorite, Colors.pinkAccent, cubit),
                    _buildActionButton(Icons.star, Colors.amber, cubit),
                    _buildActionButton(Icons.home, Colors.cyan, cubit),
                  ],
                ),
              ],
            ),
          )
        ],
      ),
    );
  }

  Widget _buildActionButton(IconData icon, Color color, IslandCubit cubit) {
    return FloatingActionButton(
      heroTag: null, // Avoid hero tag conflicts
      onPressed: () {
        // Trigger interactions via Cubit in the future
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
