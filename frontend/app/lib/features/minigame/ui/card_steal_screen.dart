// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:qlove/core/theme/app_theme.dart';

import '../bloc/card_steal_bloc.dart';

class CardStealScreen extends StatefulWidget {
  final String defenderId;
  final String targetCardId;

  const CardStealScreen({
    super.key,
    required this.defenderId,
    required this.targetCardId,
  });

  @override
  State<CardStealScreen> createState() => _CardStealScreenState();
}

class _CardStealScreenState extends State<CardStealScreen> with SingleTickerProviderStateMixin {
  late AnimationController _cardAnimController;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    _cardAnimController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 100),
      lowerBound: 0.95,
      upperBound: 1.05,
    );
    _scaleAnimation = Tween<double>(begin: 1.0, end: 0.9).animate(_cardAnimController);
    
    context.read<CardStealBloc>().add(StartStealGame(widget.defenderId, widget.targetCardId));
  }

  @override
  void dispose() {
    _cardAnimController.dispose();
    super.dispose();
  }

  void _onCardTap() {
    HapticFeedback.lightImpact();
    _cardAnimController.forward().then((_) => _cardAnimController.reverse());
    context.read<CardStealBloc>().add(TapCard());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D0D12),
      body: BlocConsumer<CardStealBloc, CardStealState>(
        listener: (context, state) {
          if (state is CardStealWon) {
            HapticFeedback.heavyImpact();
            _showResultDialog(context, true);
          } else if (state is CardStealLost) {
            HapticFeedback.heavyImpact();
            _showResultDialog(context, false);
          } else if (state is CardStealError) {
            ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(state.message)));
          }
        },
        buildWhen: (previous, current) => current is! CardStealWon && current is! CardStealLost,
        builder: (context, state) {
          if (state is CardStealLoading || state is CardStealInitial) {
            return const Center(child: CircularProgressIndicator(color: Colors.pinkAccent));
          }

          if (state is CardStealError) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.error_outline, color: Theme.of(context).colorScheme.error, size: 64),
                  const SizedBox(height: 16),
                  Text(
                    state.message,
                    style: const TextStyle(color: Colors.white, fontSize: 16),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton(
                    onPressed: () => Navigator.of(context).pop(),
                    style: ElevatedButton.styleFrom(backgroundColor: Colors.white24),
                    child: const Text('Go Back', style: TextStyle(color: Colors.white)),
                  ),
                ],
              ),
            );
          }

          if (state is CardStealPlaying) {
            final progress = state.taps / CardStealBloc.requiredTaps;
            return Stack(
              children: [
                // Background Gradient Glow
                Positioned.fill(
                  child: Container(
                    decoration: BoxDecoration(
                      gradient: RadialGradient(
                        colors: [
                          Colors.purpleAccent.withOpacity(0.2 + progress * 0.3),
                          Colors.transparent,
                        ],
                        radius: 1.2,
                      ),
                    ),
                  ),
                ),
                
                SafeArea(
                  child: Column(
                    children: [
                      const SizedBox(height: 40),
                      // Timer
                      Text(
                        '00:${state.timeLeft.toString().padLeft(2, '0')}',
                        style: GoogleFonts.inter(
                          fontSize: 48,
                          fontWeight: FontWeight.bold,
                          color: state.timeLeft <= 3 ? AppTheme.neonRed : Colors.white,
                          shadows: [
                            Shadow(
                              color: (state.timeLeft <= 3 ? AppTheme.neonRed : AppTheme.neonBlue).withOpacity(0.8),
                              blurRadius: 15,
                            )
                          ],
                        ),
                      ),
                      const SizedBox(height: 20),
                      Text(
                        'TAP TO STEAL!',
                        style: GoogleFonts.outfit(
                          fontSize: 24,
                          fontWeight: FontWeight.w600,
                          color: Colors.white70,
                          letterSpacing: 2,
                        ),
                      ),
                      
                      const Spacer(),
                      
                      // The Card
                      GestureDetector(
                        onTap: _onCardTap,
                        child: ScaleTransition(
                          scale: _scaleAnimation,
                          child: Container(
                            width: 260,
                            height: 380,
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(24),
                              border: Border.all(
                                color: Colors.white.withOpacity(0.2),
                                width: 2,
                              ),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.pinkAccent.withOpacity(0.3 + progress * 0.5),
                                  blurRadius: 30,
                                  spreadRadius: progress * 20,
                                ),
                              ],
                              gradient: const LinearGradient(
                                begin: Alignment.topLeft,
                                end: Alignment.bottomRight,
                                colors: [
                                  Color(0xFF2B2B36),
                                  Color(0xFF1E1E26),
                                ],
                              ),
                            ),
                            child: ClipRRect(
                              borderRadius: BorderRadius.circular(24),
                              child: BackdropFilter(
                                filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                                child: Center(
                                  child: Icon(
                                    Icons.diamond_outlined,
                                    size: 100,
                                    color: Colors.white.withOpacity(0.8),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                      
                      const Spacer(),
                      
                      // Power Bar
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 40),
                        child: Column(
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text('Power', style: GoogleFonts.inter(color: Colors.white70)),
                                Text('${(progress * 100).toInt()}%', style: GoogleFonts.inter(color: Colors.white, fontWeight: FontWeight.bold)),
                              ],
                            ),
                            const SizedBox(height: 8),
                            ClipRRect(
                              borderRadius: BorderRadius.circular(10),
                              child: LinearProgressIndicator(
                                value: progress,
                                minHeight: 12,
                                backgroundColor: Colors.white12,
                                valueColor: const AlwaysStoppedAnimation<Color>(Colors.pinkAccent),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 40),
                    ],
                  ),
                ),
              ],
            );
          }
          return const SizedBox();
        },
      ),
    );
  }

  void _showResultDialog(BuildContext context, bool isWin) {
    showGeneralDialog(
      context: context,
      barrierDismissible: false,
      transitionDuration: const Duration(milliseconds: 500),
      pageBuilder: (context, anim1, anim2) {
        return Scaffold(
          backgroundColor: Colors.black87,
          body: Center(
            child: TweenAnimationBuilder<double>(
              tween: Tween(begin: 0.0, end: 1.0),
              duration: const Duration(milliseconds: 800),
              curve: Curves.elasticOut,
              builder: (context, value, child) {
                return Transform.scale(
                  scale: value,
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isWin ? Icons.check_circle_outline : Icons.cancel_outlined,
                        size: 120,
                        color: isWin ? AppTheme.neonGreen : AppTheme.neonRed,
                      ),
                      const SizedBox(height: 24),
                      Text(
                        isWin ? 'SUCCESS!' : 'FAILED!',
                        style: GoogleFonts.outfit(
                          fontSize: 48,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                          shadows: [
                            Shadow(
                              color: (isWin ? AppTheme.neonGreen : AppTheme.neonRed).withOpacity(0.6),
                              blurRadius: 20,
                            )
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                      Text(
                        isWin ? 'You have stolen the card.' : '+500 Xu Consolation',
                        style: GoogleFonts.inter(
                          fontSize: 20,
                          color: Colors.white70,
                        ),
                      ),
                      const SizedBox(height: 40),
                      ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.white24,
                          padding: const EdgeInsets.symmetric(horizontal: 40, vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(30),
                          ),
                        ),
                        onPressed: () {
                          Navigator.of(context).pop();
                          Navigator.of(context).pop(); // Go back to previous screen
                        },
                        child: Text(
                          'CONTINUE',
                          style: GoogleFonts.inter(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
          ),
        );
      },
    );
  }
}
